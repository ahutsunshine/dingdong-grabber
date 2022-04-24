package strategy

import (
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

func (sf *schedulerFactory) Build(strategy int, u *user.User, baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis int,
	crons []string, play bool, pushToken string) Interface {
	switch strategy {
	case 0: // 人工策略
		return NewManualScheduler(order.NewOrder(u, order.ManualStrategy), baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis, play, pushToken)
	case 1: // 定时策略
		return NewTimingScheduler(order.NewOrder(u, order.TimingStrategy), baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis, crons, play, pushToken)
	default:
		klog.Fatalf("不支持此无效策略: %d", strategy)
	}
	return nil
}
