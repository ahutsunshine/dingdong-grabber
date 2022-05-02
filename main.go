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
	"flag"

	schedule "github.com/dingdong-grabber/pkg/strategy"
	"github.com/dingdong-grabber/pkg/user"
	"k8s.io/klog"
)

const (
	defaultBaseThreadSize        = 2   // 默认基础信息执行线程数, 建议不要超过2，否则容易被风控
	defaultSubmitOrderThreadSize = 2   // 默认提交订单执行线程数, 建议不要超过2，否则容易被风控
	defaultMinSleepMillis        = 300 // 默认抢菜最小时间间隔 300ms
	defaultMaxSleepMillis        = 500 // 默认抢菜最大时间间隔 500ms

	// 抢菜策略, 0: 人工策略，1: 定时策略, 默认是定时策略，2: 哨兵策略
	// - 人工策略: 程序运行即开始抢菜, 此策略下程序默认只会跑2分钟，如果没有商品库存，则会立即停止
	// - 定时策略: 定时抢菜，事先订好时间，叮咚默认是早上5:59:50和8:29:50开始抢菜，这种模式要避免启动过早导致用户登录信息过期。
	// - 哨兵策略: 捡漏模式，长期运行捡漏可配送时间
	// 使用cron job定义抢菜时间
	// 每天 5:59:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 59 05  * *  *

	// 每天 08:29:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 29 08  * *  *
	strategy = 1

	// 必须填写用户cookie， cookie代表人的身份
	cookie = "" // 请求头部的Cookie

	// 可选参数，可不填。pushplus token用于推送下单成功消息, http://www.pushplus.plus
	pushToken = ""
)

var (
	c    string
	sy   int
	play bool
)

// 更新兼容叮咚微信小程序版本：2.85.2

func main() {
	flag.StringVar(&c, "cookie", "", "请求头部的Cookie")
	flag.IntVar(&sy, "strategy", 1, "设置抢菜策略")
	// 抢菜成功后是否播放《Everything I Need》通知用户
	// false: 不播放
	// true: 播放
	flag.BoolVar(&play, "play", true, "抢菜成功后播放音乐通知用户")
	flag.Parse()

	setDefault()

	// 1. 初始化用户必须的参数数据
	u := user.NewDefaultUser()
	if err := u.LoadConfig(c); err != nil {
		return
	}

	// 2. 构建实际调度策略
	factory := schedule.NewSchedulerFactory()
	scheduler := factory.Build(strategy, u, defaultBaseThreadSize, defaultSubmitOrderThreadSize,
		defaultMinSleepMillis, defaultMaxSleepMillis, []string{"30 59 05 * * ?", "30 29 08 * * ?"}, play, pushToken)

	// 3. 运行调度策略抢菜
	if err := scheduler.Schedule(context.TODO()); err != nil {
		klog.Fatal(err)
	}

	select {}
}

// setDefault 是为了方便用户在main上方直接填写用户必要参数
func setDefault() {
	if c == "" {
		c = cookie
	}
}
