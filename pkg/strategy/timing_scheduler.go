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
	"strings"

	"github.com/dingdong-grabber/pkg/order"
	"github.com/robfig/cron/v3"
	"k8s.io/klog"
)

// TimingScheduler 定时策略调度器
type TimingScheduler struct {
	Scheduler `json:",inline"`
	crons     []string // cron job 调度时间
}

func NewTimingScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis int, crons []string, play bool, pushToken string) Interface {
	if minSleepMillis > maxSleepMillis {
		maxSleepMillis = minSleepMillis
	}
	return &TimingScheduler{
		Scheduler: Scheduler{
			o:                    o,
			play:                 play,
			baseTheadSize:        baseTheadSize,
			submitOrderTheadSize: submitOrderTheadSize,
			minSleepMillis:       minSleepMillis,
			maxSleepMillis:       maxSleepMillis,
			pushToken:            pushToken,
		},
		crons: crons,
	}
}

// Schedule 使用cron调度
func (ts *TimingScheduler) Schedule(ctx context.Context) error {
	c := cron.New(cron.WithSeconds())

	// 定时任务每隔3s需要检测token是否过期，过期则直接退出
	if _, err := c.AddFunc("0/3 * * * * *", func() {
		if _, err := ts.Scheduler.o.User().GetDefaultAddr(); err != nil && strings.Contains(err.Error(), "已过期") {
			klog.Fatal("用户Cookie已过期，请重新填写")
		}
	}); err != nil {
		klog.Error(err)
		return err
	}

	// 定义的定时任务
	for _, spec := range ts.crons {
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
