package strategy

import "context"

type Interface interface {
	Schedule(ctx context.Context) error
}
