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

package main

import (
	"context"

	"github.com/dingdong-grabber/pkg/config"
	schedule "github.com/dingdong-grabber/pkg/strategy"
	"github.com/dingdong-grabber/pkg/user"
	"k8s.io/klog"
)

// 更新兼容叮咚微信小程序版本：2.85.2
func main() {
	// 1. 加载全局配置: 根目录config.yaml，用户只需配置cookie参数，其余均是可循参数
	m := config.Manager{}
	if err := m.LoadConfig(); err != nil {
		klog.Fatal()
	}

	// 2. 初始化用户必须的参数数据
	u := user.NewDefaultUser()
	if err := u.LoadConfig(m.Conf.Config.Cookie); err != nil {
		return
	}

	// 3. 构建实际调度策略
	factory := schedule.NewSchedulerFactory()
	scheduler := factory.Build(m.Conf.Config, u)

	// 4. 运行调度策略抢菜
	if err := scheduler.Schedule(context.TODO()); err != nil {
		klog.Fatal(err)
	}

	select {}
}
