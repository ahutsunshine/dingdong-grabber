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
	"github.com/dingdong-grabber/pkg/config"
	"github.com/dingdong-grabber/pkg/order"
	"github.com/dingdong-grabber/pkg/user"
	"k8s.io/klog"
)

type schedulerFactory struct {
}

type SchedulerFactory interface {
	Build() Interface
}

func NewSchedulerFactory() *schedulerFactory {
	return &schedulerFactory{}
}

func (sf *schedulerFactory) Build(c *config.Config, u *user.User) Interface {
	switch c.Strategy {
	case 0: // 人工策略
		return NewManualScheduler(order.NewOrder(u, order.ManualStrategy), c)
	case 1: // 定时策略
		return NewTimingScheduler(order.NewOrder(u, order.TimingStrategy), c)
	case 2: // 哨兵策略
		return NewSentinelScheduler(order.NewOrder(u, order.SentinelStrategy), c)
	default:
		klog.Fatalf("不支持此无效策略: %d", c.Strategy)
	}
	return nil
}
