package downlink

import (
	"context"
	"time"
)

type Backoff struct{ Delay time.Duration }

func (b Backoff) Wait(ctx context.Context) error {
	timer := time.NewTimer(b.Delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
