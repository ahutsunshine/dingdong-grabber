package strategy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dingdong-grabber/pkg/notice"
	"github.com/dingdong-grabber/pkg/order"
	"k8s.io/klog"
)

type Scheduler struct {
	o                    *order.Order
	play                 bool    // 播放音乐按钮
	minOrderPrice        float64 // 最小订单成交金额
	baseTheadSize        int     // 基础信息执行线程数
	submitOrderTheadSize int     // 提交订单执行线程数
	minSleepMillis       int     // 请求间隔时间最小值
	maxSleepMillis       int     // 请求间隔时间最大值
	pushToken            string
}

// Run 作为保护线程负责检查订单是否下单成功，2分钟未下单自动终止,避免对叮咚服务器造成压力,也避免封号
func (s *Scheduler) Run(ctx context.Context) {
	go func() {
		var (
			deadline = time.After(120 * time.Second)
			ticker   = time.NewTicker(3 * time.Second)
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
				if _, err := s.o.User().GetUserDetail(); err != nil && strings.Contains(err.Error(), "已过期") {
					klog.Fatal("用户Cookie已过期，请重新填写")
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
				// 购物车选中后，可能以后并不需要再次选中(但不确定背后逻辑，再每隔1-3秒选中一次)
				time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
			}
		}()

		go func() {
			for !s.o.Stop() {
				time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
				cart, err := s.o.GetCart()
				if err != nil {
					continue
				}
				// 购物车无可购买的商品
				if cart == nil {
					continue
				}
				if cart["total_money"] == nil {
					bytes, err := json.Marshal(cart)
					if err != nil {
						klog.Errorf("解析购物车信息出错, 错误: %v", err)
					} else {
						klog.Infof("获取购物总金额出错，购物车无总金额参数, 详情: %s", string(bytes))
					}
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
					klog.Infof("订单金额：%f, 不满足最小金额设置：%f, 继续重试", cart["total_money"], s.minOrderPrice)
				} else {
					s.o.SetCart(cart)
				}
			}
		}()

		go func() {
			for !s.o.Stop() {
				time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
				if s.o.Cart() == nil {
					continue
				}
				reservedTimes, err := s.o.GetMultiReserveTime()
				if err != nil {
					continue
				}
				s.o.SetReservedTime(reservedTimes)
			}
		}()

		go func() {
			for !s.o.Stop() {
				time.Sleep(time.Duration(rand.Intn(s.maxSleepMillis-s.minSleepMillis)+s.minSleepMillis) * time.Millisecond)
				if s.o.Cart() == nil || s.o.ReservedTime() == nil {
					continue
				}
				checkOrder, err := s.o.GetCheckOrder()
				if err != nil {
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
					continue
				}
				_, err := s.o.SubmitOrder()
				if err != nil {
					continue
				}
				// 下单已成功，停止抢菜
				s.o.SetStop(true)
				klog.Infof("下单成功，请在5分钟内支付金额: %s，否则订单会被叮咚自动取消", s.o.Cart()["total_money"])

				go func() {
					// 播放音乐通知用户
					if s.play {
						mp3 := notice.NewDefaultMp3()
						if err = mp3.Notify(); err != nil {
							klog.Error(err)
						}
					}
				}()

				if s.pushToken != "" {
					// 推送一次即可，失败则重试2次
					for i := 0; i < 3; i++ {
						p := notice.NewPush(s.pushToken, "抢菜已成功，请前往APP付款", fmt.Sprintf("下单成功，请在5分钟内支付金额: %v，否则订单会被叮咚自动取消", s.o.Cart()["total_money"]))
						if err := p.Notify(); err == nil {
							break
						}
					}
				}

				// 休眠30s, 让音乐飞一会
				time.Sleep(time.Second * 30)

				// 正常退出程序
				os.Exit(0)
			}
		}()
	}
	return nil
}
