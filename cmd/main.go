package main

import (
	"context"
	"github.com/dingdong-grabber/pkg/order"
	"github.com/dingdong-grabber/pkg/strategy"

	"github.com/dingdong-grabber/pkg/user"
)

const (
	defaultBaseThreadSize        = 2 // 默认基础信息执行线程数
	defaultSubmitOrderThreadSize = 4 // 默认提交订单执行线程数
)

// 当前小程序版本：2.83.0

func main() {
	u := user.NewDefaultUser()
	// 1. 初始化用户必须的参数数据
	if err := u.LoadConfig(); err != nil {
		return
	}

	// 人工模式
	//o := order.NewOrder(u, order.ManualStrategy)
	//scheduler := strategy.NewManualScheduler(o, defaultBaseThreadSize, defaultSubmitOrderThreadSize, 1000, 2000)

	// 定时模式
	// cron job定义
	// 每天 5:59:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 59 05  * *  *

	// 每天 08:29:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 29 08  * *  *
	o := order.NewOrder(u, order.ManualStrategy)
	scheduler := strategy.NewTimingScheduler(o, defaultBaseThreadSize, defaultSubmitOrderThreadSize, 2000, 3000, []string{"01 51 17 * * ?"})

	_ = scheduler.Schedule(context.TODO())

	select {}
}
