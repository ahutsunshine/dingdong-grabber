package strategy

import (
	"context"

	"github.com/dingdong-grabber/pkg/order"
)

// ManualScheduler 人工策略调度器
type ManualScheduler struct {
	Scheduler `json:",inline"`
}

func NewManualScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis int, play bool, pushToken string) Interface {
	if minSleepMillis > maxSleepMillis {
		maxSleepMillis = minSleepMillis
	}
	return &ManualScheduler{Scheduler{
		o:                    o,
		play:                 play,
		baseTheadSize:        baseTheadSize,
		submitOrderTheadSize: submitOrderTheadSize,
		minSleepMillis:       minSleepMillis,
		maxSleepMillis:       maxSleepMillis,
		pushToken:            pushToken,
	}}
}

func (ms *ManualScheduler) Schedule(ctx context.Context) error {
	return ms.Scheduler.Schedule(ctx)
}
