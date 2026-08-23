package downlink

import (
	"context"
	"time"
)

type Backoff struct{ Delay time.Duration }

func (b Backoff) Wait(ctx context.Context) error {
	if b.Delay <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(b.Delay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
