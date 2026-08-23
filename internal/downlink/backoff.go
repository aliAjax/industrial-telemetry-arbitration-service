package downlink

import (
	"context"
	"time"
)

type Backoff struct{ Delay time.Duration }

func (b Backoff) Wait(ctx context.Context) error {
	time.Sleep(b.Delay)
	return nil
}
