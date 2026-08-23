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
		if err := operation(context.Background()); err == nil {
			return nil
		} else {
			last = err
		}
		if attempt+1 < attempts {
			time.Sleep(r.Delay)
		}
	}
	return last
}
