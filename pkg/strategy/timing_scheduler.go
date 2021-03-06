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

	"github.com/dingdong-grabber/pkg/config"
	"github.com/dingdong-grabber/pkg/order"
	"github.com/robfig/cron/v3"
	"k8s.io/klog"
)

// TimingScheduler 定时策略调度器
type TimingScheduler struct {
	Scheduler `json:",inline"`
	cronJobs  []string // cron job 调度时间
}

func NewTimingScheduler(o *order.Order, c *config.Config) Interface {
	var (
		minSleepMillis = c.MinSleepMillis
		maxSleepMillis = c.MaxSleepMillis
	)
	if minSleepMillis > c.MaxSleepMillis {
		maxSleepMillis = minSleepMillis
	}
	return &TimingScheduler{
		Scheduler: Scheduler{
			o:                    o,
			play:                 c.Play,
			baseTheadSize:        c.BaseThreadSize,
			submitOrderTheadSize: c.SubmitOrderThreadSize,
			minSleepMillis:       minSleepMillis,
			maxSleepMillis:       maxSleepMillis,
			pushToken:            c.PushToken,
		},
		cronJobs: c.CronJobs,
	}
}

// Schedule 使用cron调度
func (ts *TimingScheduler) Schedule(ctx context.Context) error {
	klog.Info("正在使用定时模式，默认在5:59:40或者08:29:40开始抢菜")
	c := cron.New(cron.WithSeconds())
	// 定义的定时任务
	for _, spec := range ts.cronJobs {
		if _, err := c.AddFunc(spec, func() {
			_ = ts.Scheduler.Schedule(ctx)
		}); err != nil {
			klog.Error(err)
			return err
		}
	}
	c.Start()

	return nil
}
