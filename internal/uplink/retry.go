package uplink

import (
	"context"
	"time"
)

type Retryer struct {
	Attempts int
	Delay    time.Duration
}

func (r Retryer) Do(ctx context.Context, operation func(context.Context) error) error {
	attempts := r.Attempts
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := operation(ctx); err == nil {
			return nil
		} else {
			last = err
		}
		if attempt+1 < attempts && r.Delay > 0 {
			timer := time.NewTimer(r.Delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return last
}
