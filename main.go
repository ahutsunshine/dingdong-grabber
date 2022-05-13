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

// 请注意，请根据README中的教程在config.yaml中填写必要配置。如有必要，请在charles/文件夹下配置*.chlsj基础文件，此文件包含了用户手机真实环境，
// 可较大可能避免被风控，但还是治标不治本，不可长时间运行。

// **************************
// 非常重要！！非常重要！！非常重要！！
// 建议抢菜前，在config.yaml中将strategy设置为3(测试模式)，用于测试配置是否正确，流程是否已通。
// 如果不通，请不要运行其他策略，避免风控。
// 1. 目前仅可以支持ios设备
// 2. 代码还有瑕疵需要完善，最近太累了，让我缓一缓
// 3. 代码此版本仅做参考，后续会继续更新
// **************************

func main() {
	// 1. 加载全局配置: 根目录config.yaml，用户只需配置cookie参数，其余均是可循参数
	m := config.Manager{}
	if err := m.LoadConfig(); err != nil {
		klog.Fatal(err)
	}

	// 2. 初始化用户必须的参数数据
	u := user.NewDefaultUser(m.Conf.Config)
	if err := u.LoadConfig(); err != nil {
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
