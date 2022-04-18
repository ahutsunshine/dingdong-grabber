package strategy

import (
	"context"

	"github.com/dingdong-grabber/pkg/order"
	"github.com/robfig/cron/v3"
)

type TimingScheduler struct {
	Scheduler `json:",inline"`
	crons     []string // cron job 调度时间
}

func NewDefaultTimingScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize int, crons []string) Interface {
	return NewTimingScheduler(o, baseTheadSize, submitOrderTheadSize, DefaultMinSleepMillis, DefaultMaxSleepMillis, crons)
}

func NewTimingScheduler(o *order.Order, baseTheadSize, submitOrderTheadSize, minSleepMillis, maxSleepMillis int, crons []string) Interface {
	if minSleepMillis > maxSleepMillis {
		maxSleepMillis = minSleepMillis
	}
	return &TimingScheduler{
		Scheduler: Scheduler{
			o:                    o,
			baseTheadSize:        baseTheadSize,
			submitOrderTheadSize: submitOrderTheadSize,
			minSleepMillis:       minSleepMillis,
			maxSleepMillis:       maxSleepMillis,
		},
		crons: crons,
	}
}

// Schedule 使用cron调度
func (ts *TimingScheduler) Schedule(ctx context.Context) error {
	c := cron.New(cron.WithSeconds())
	for _, spec := range ts.crons {
		if _, err := c.AddFunc(spec, func() {
			ts.Scheduler.Schedule(ctx)
		}); err != nil {
			return err
		}
	}
	c.Start()
	return nil
}
