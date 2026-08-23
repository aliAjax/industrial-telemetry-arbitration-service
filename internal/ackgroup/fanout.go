package ackgroup

import (
	"context"
	"sync"
)

type Processor func(context.Context, Ack) (Ack, error)

func Fanout(ctx context.Context, inputs []Ack, processor Processor) (<-chan Ack, <-chan error) {
	results := make(chan Ack, len(inputs))
	errorsOut := make(chan error, len(inputs))
	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, input := range inputs {
		input := input.Clone()
		go func() {
			defer wg.Done()
			result, err := processor(ctx, input)
			if err != nil {
				errorsOut <- err
				return
			}
			select {
			case results <- result.Clone():
			case <-ctx.Done():
				errorsOut <- ctx.Err()
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
		close(errorsOut)
	}()
	return results, errorsOut
}
