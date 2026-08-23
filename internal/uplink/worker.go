package uplink

import "context"

type Worker struct {
	retry  Retryer
	sender Sender
}

func NewWorker(retry Retryer, sender Sender) *Worker { return &Worker{retry: retry, sender: sender} }

func (w *Worker) Process(ctx context.Context, requests []Request) error {
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return err
		}
		current := request
		if err := w.retry.Do(ctx, func(callCtx context.Context) error {
			return Handle(callCtx, w.sender, current)
		}); err != nil {
			return err
		}
	}
	return nil
}
