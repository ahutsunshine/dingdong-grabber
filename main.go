package main

import (
	"context"
	"flag"
	schedule "github.com/dingdong-grabber/pkg/strategy"
	"github.com/dingdong-grabber/pkg/user"
)

const (
	defaultBaseThreadSize        = 2 // 默认基础信息执行线程数
	defaultSubmitOrderThreadSize = 4 // 默认提交订单执行线程数

	// 抢菜策略, 0: 人工策略，1: 定时策略, 默认是定时策略
	// - 人工策略: 程序运行即开始抢菜, 此策略下程序默认只会跑2分钟，如果没有商品库存，则会立即停止
	// - 定时策略: 定时抢菜，事先订好时间，叮咚默认是早上5:59:50和8:29:50开始抢菜，这种模式要避免启动过早导致用户登录信息过期。
	// 使用cron job定义抢菜时间
	// 每天 5:59:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 59 05  * *  *

	// 每天 08:29:50秒的时候开始执行
	// 秒 分  时 日 月 周
	// 50 29 08  * *  *
	strategy = 1

	deviceId    = ""
	cookie      = ""
	uid         = ""
	userAgent   = ""
	sid         = ""
	deviceToken = ""
)

var (
	di, c, u, ua, s, dt string
	sy                  int
)

// 应用小程序版本：2.83.0
func main() {
	flag.StringVar(&di, "device-id", "", "请求头部的ddmc-device-id")
	flag.StringVar(&c, "cookie", "", "请求头部的Cookie")
	flag.StringVar(&u, "uid", "", "请求头部的ddmc-id")
	flag.StringVar(&ua, "user-agent", "", "请求头部的User-Agent")
	flag.StringVar(&s, "sid", "", "请求信息的s_id")
	flag.StringVar(&dt, "device-token", "", "请求信息的device_token")
	flag.IntVar(&sy, "strategy", 1, "设置抢菜策略")
	setDefault()

	// 1. 初始化用户必须的参数数据
	user := user.NewDefaultUser()
	if err := user.LoadConfig(di, c, u, ua, s, dt); err != nil {
		return
	}

	// 2. 构建实际调度策略
	factory := schedule.NewSchedulerFactory()
	scheduler := factory.Build(strategy, user, defaultBaseThreadSize, defaultSubmitOrderThreadSize, 300, 500, []string{"50 59 05 * * ?", "50 29 08 * * ?"})

	// 3. 运行调度策略抢菜
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
