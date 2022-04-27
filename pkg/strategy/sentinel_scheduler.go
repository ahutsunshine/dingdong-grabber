/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package strategy

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dingdong-grabber/pkg/notice"
	"github.com/dingdong-grabber/pkg/order"
	"k8s.io/klog"
)

const defaultRepeatCount = 5

// SentinelScheduler 捡漏策略调度器
type SentinelScheduler struct {
	Scheduler   `json:",inline"`
	repeatCount int
}

func NewSentinelScheduler(o *order.Order, minSleepMillis, maxSleepMillis int, play bool, pushToken string) Interface {
	if minSleepMillis < 5000 {
		minSleepMillis = 5000
		klog.Info("使用默认哨兵策略，每隔5-15s发起一起请求")
	}
	if maxSleepMillis <= minSleepMillis {
		maxSleepMillis = minSleepMillis + 10000
	}

	return &SentinelScheduler{
		Scheduler: Scheduler{
			o:              o,
			play:           play,
			minSleepMillis: minSleepMillis,
			maxSleepMillis: maxSleepMillis,
			pushToken:      pushToken,
		},
		repeatCount: defaultRepeatCount,
	}
}

func (ss *SentinelScheduler) Schedule(ctx context.Context) error {
	var loopCount int
	for !ss.o.Stop() {
		time.Sleep(time.Duration(rand.Intn(ss.maxSleepMillis-ss.minSleepMillis)+ss.minSleepMillis) * time.Millisecond)
		loopCount++
		// 每循环抢菜60次就休会1-3分钟
		if loopCount%60 == 0 {
			time.Sleep(time.Duration(rand.Intn(60000)+2*60000) * time.Millisecond)
		}

		var err error
		for i := 0; i < ss.repeatCount; i++ {
			if err = ss.o.CheckAll(); err == nil {
				break
			} else if err != nil {
				//随机等待100-300毫秒
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
				continue
			}
			break
		}
		// 只要上一步出错就重新开始
		if err != nil {
			continue
		}

		for i := 0; i < ss.repeatCount; i++ {
			var cart map[string]interface{}
			if cart, err = ss.o.GetCart(); err == nil {
				// 购物车无可购买的商品
				if cart == nil {
					err = errors.New("购物车无可购买的商品")
					break
				}
				ss.o.SetCart(cart)
				break
			} else if err != nil {
				//随机等待100-300毫秒
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			continue
		}

		for i := 0; i < ss.repeatCount; i++ {
			var reservedTimes map[string]interface{}
			if reservedTimes, err = ss.o.GetMultiReserveTime(); err == nil {
				if reservedTimes == nil {
					err = errors.New("当前运力紧张，今天各时段已约满")
				}
				ss.o.SetReservedTime(reservedTimes)
				break
			} else if err != nil {
				//随机等待100-300毫秒
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			continue
		}

		for i := 0; i < ss.repeatCount; i++ {
			var checkOrder map[string]interface{}
			if checkOrder, err = ss.o.GetCheckOrder(); err == nil {
				ss.o.SetCheckOrder(checkOrder)
				break
			} else if err != nil {
				//随机等待100-300毫秒
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			continue
		}

		for i := 0; i < ss.repeatCount; i++ {
			if _, err = ss.o.SubmitOrder(); err == nil {
				// 下单已成功，停止抢菜
				ss.o.SetStop(true)
				break
			} else if err != nil {
				//随机等待100-300毫秒
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			continue
		}

		klog.Infof("下单成功，请在5分钟内支付金额: %s，否则订单会被叮咚自动取消", ss.o.Cart()["total_money"])

		go func() {
			// 播放音乐通知用户
			if ss.play {
				mp3 := notice.NewDefaultMp3()
				if err = mp3.Notify(); err != nil {
					klog.Error(err)
				}
			}
		}()

		if ss.pushToken != "" {
			// 推送一次即可，失败则重试2次
			for i := 0; i < 3; i++ {
				p := notice.NewPush(ss.pushToken, "抢菜已成功，请前往APP付款", fmt.Sprintf("下单成功，请在5分钟内支付金额: %v，否则订单会被叮咚自动取消", ss.o.Cart()["total_money"]))
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
	return nil
}
