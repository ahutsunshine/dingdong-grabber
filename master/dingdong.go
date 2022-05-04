package main

import (
	"context"
	"github.com/dingdong-grabber/pkg/config"
	schedule "github.com/dingdong-grabber/pkg/strategy"
	"github.com/dingdong-grabber/pkg/user"
	"k8s.io/klog"
	"sync"
)

// 更新兼容叮咚微信小程序版本：2.85.2
func ddMain(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	defer waitGroup.Done()

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
