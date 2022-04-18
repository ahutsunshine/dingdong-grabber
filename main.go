package main

import (
	"context"
	"flag"
	"github.com/dingdong-grabber/pkg/order"
	"github.com/dingdong-grabber/pkg/strategy"

	"github.com/dingdong-grabber/pkg/user"
)

const (
	defaultBaseThreadSize        = 2 // 默认基础信息执行线程数
	defaultSubmitOrderThreadSize = 4 // 默认提交订单执行线程数

	// -------Header 请求头部必填项----------
	deviceId  = "" // 请求头部的ddmc-device-id
	cookie    = "" // 请求头部的Cookie
	uid       = "" // 请求头部的ddmc-uid
	userAgent = "" // 请求头部的User-Agent

	// -------Body 请求信息必填项----------
	sid         = "" // 请求信息的s_id
	deviceToken = "" // 请求信息的device_token
)

var di, c, u, ua, s, dt string

// 当前应用小程序版本：2.83.0
func main() {
	flag.StringVar(&di, "device-id", "", "请求头部的ddmc-device-id")
	flag.StringVar(&c, "cookie", "", "请求头部的Cookie")
	flag.StringVar(&u, "uid", "", "请求头部的ddmc-id")
	flag.StringVar(&ua, "user-agent", "", "请求头部的User-Agent")
	flag.StringVar(&s, "sid", "", "请求信息的s_id")
	flag.StringVar(&dt, "device-token", "", "请求信息的device_token")
	setDefault()

	user := user.NewDefaultUser()
	// 1. 初始化用户必须的参数数据
	if err := user.LoadConfig(di, c, u, ua, s, dt); err != nil {
		return
	}

	// 2.1 人工模式
	//o := order.NewOrder(u, order.ManualStrategy)
	//scheduler := strategy.NewManualScheduler(o, defaultBaseThreadSize, defaultSubmitOrderThreadSize, 1000, 2000)

	// 2.2 定时模式
	// cron job定义
	// 每天 5:59:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 59 05  * *  *

	// 每天 08:29:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 29 08  * *  *
	o := order.NewOrder(user, order.TimingStrategy)
	scheduler := strategy.NewTimingScheduler(o, defaultBaseThreadSize, defaultSubmitOrderThreadSize, 300, 500, []string{"50 59 05 * * ?", "50 29 08 * * ?"})

	_ = scheduler.Schedule(context.TODO())

	select {}
}

// setDefault 是为了方便用户在main上方直接填写用户必要参数
func setDefault() {
	if di == "" {
		di = deviceId
	}
	if c == "" {
		c = cookie
	}
	if u == "" {
		u = uid
	}
	if ua == "" {
		ua = userAgent
	}
	if s == "" {
		s = sid
	}
	if dt == "" {
		dt = deviceToken
	}
}
