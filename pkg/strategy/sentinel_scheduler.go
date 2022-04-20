package strategy

import "context"

// SentinelScheduler 捡漏策略调度器
type SentinelScheduler struct {
	Scheduler `json:",inline"`
}

func (ss *SentinelScheduler) Schedule(ctx context.Context) error {
	return nil
}
