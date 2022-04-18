package strategy

import (
	"context"

	"github.com/dingdong-grabber/pkg/order"
)

type ManualScheduler struct {
	Scheduler `json:",inline"`
}

func NewDefaultManualScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize int) Interface {
	return NewManualScheduler(o, baseTheadSize, submitOrderTheadSize, DefaultMinSleepMillis, DefaultMaxSleepMillis)
}

func NewManualScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis int) Interface {
	if minSleepMillis > maxSleepMillis {
		maxSleepMillis = minSleepMillis
	}
	return &ManualScheduler{Scheduler{
		o:                    o,
		baseTheadSize:        baseTheadSize,
		submitOrderTheadSize: submitOrderTheadSize,
		minSleepMillis:       minSleepMillis,
		maxSleepMillis:       maxSleepMillis,
	}}
}

func (ms *ManualScheduler) Schedule(ctx context.Context) error {
	return ms.Scheduler.Schedule(ctx)
}
