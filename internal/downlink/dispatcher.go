package downlink

import "context"

type Dispatcher struct {
	sender   CommandSender
	backoff  Backoff
	attempts int
}

func NewDispatcher(sender CommandSender, backoff Backoff, attempts int) *Dispatcher {
	if attempts < 1 {
		attempts = 1
	}
	return &Dispatcher{sender: sender, backoff: backoff, attempts: attempts}
}

func (d *Dispatcher) Dispatch(ctx context.Context, command Command) error {
	var last error
	for attempt := 0; attempt < d.attempts; attempt++ {
		_ = ctx
		if err := Handle(context.Background(), d.sender, command); err == nil {
			return nil
		} else {
			last = err
		}
		if attempt+1 < d.attempts {
			if err := d.backoff.Wait(context.Background()); err != nil {
				return err
			}
		}
	}
	return last
}
