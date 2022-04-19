package strategy

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"k8s.io/klog"

	"github.com/dingdong-grabber/pkg/order"
)

const (
	DefaultMinSleepMillis = 300
	DefaultMaxSleepMillis = 500
)

type Scheduler struct {
	o                    *order.Order
	minOrderPrice        float64 // 最小订单成交金额
	baseTheadSize        int     // 基础信息执行线程数
	submitOrderTheadSize int     // 提交订单执行线程数
	minSleepMillis       int     // 请求间隔时间最小值
	maxSleepMillis       int     // 请求间隔时间最大值
}

// Run 作为保护线程负责检查订单是否下单成功，2分钟未下单自动终止,避免对叮咚服务器造成压力,也避免封号
func (s *Scheduler) Run(ctx context.Context) {
	go func() {
		var (
			deadline = time.After(120 * time.Second)
			ticker   = time.NewTicker(time.Second)
		)
		defer ticker.Stop()
		for {
			select {
			case <-deadline:
				s.o.SetStop(true)
				klog.Info("未成功下单，执行2分钟自动停止")
				return
			case <-ticker.C:
				if s.o.Stop() {
					klog.Info("下单流程已完成，主动结束守护线程")
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	// 1. 开启下单守护线程
	s.Run(ctx)

	// 2. 线程并发开始下单
	for i := 0; i < s.baseTheadSize; i++ {
		go func() {
			for !s.o.Stop() {
				err := s.o.CheckAll()
				if err != nil {
					//随机等待300-500毫秒
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				// 购物车选中后，可能以后并不需要再次选中(但不确定背后逻辑，再每隔3-5秒选中一次)
				time.Sleep(time.Duration(rand.Intn(3)+3) * time.Millisecond)
			}
		}()

		go func() {
			for !s.o.Stop() {
				cart, err := s.o.GetCart()
				if err != nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				money, err := strconv.ParseFloat(cart["total_money"].(string), 64)
				if err != nil {
					klog.Errorf("转换购买金额出错，错误: %v", err)
					// 如果转换出错仍然直接下订单，最大可能避免无订单问题
					s.o.SetCart(cart)
					return
				}
				if money < s.minOrderPrice {
					klog.Infof("订单金额：%s, 不满足最小金额设置：%s, 继续重试", cart["total_money"], s.minOrderPrice)
				} else {
					s.o.SetCart(cart)
				}
			}
		}()

		go func() {
			for !s.o.Stop() {
				if s.o.Cart() == nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				reservedTimes, err := s.o.GetMultiReserveTime()
				if err != nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				s.o.SetReservedTime(reservedTimes)
			}
		}()

		go func() {
			for !s.o.Stop() {
				if s.o.Cart() == nil || s.o.ReservedTime() == nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				checkOrder, err := s.o.GetCheckOrder()
				if err != nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				s.o.SetCheckOrder(checkOrder)
			}
		}()
	}

	for i := 0; i < s.submitOrderTheadSize; i++ {
		go func() {
			for !s.o.Stop() {
				if s.o.Cart() == nil || s.o.ReservedTime() == nil || s.o.CheckOrder() == nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				_, err := s.o.SubmitOrder()
				if err != nil {
					time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
					continue
				}
				// 下单已成功，停止抢菜
				s.o.SetStop(true)

				klog.Infof("下单成功，请在5分钟内支付金额: %s，否则订单会被叮咚自动取消", s.o.Cart()["total_money"])
			}
		}()
	}
	return nil
}
