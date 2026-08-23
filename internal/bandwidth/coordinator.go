package bandwidth

import "sync"

type Coordinator struct{}

func (Coordinator) Rebalance(producers []Producer) ([]<-chan Allocation, <-chan error, <-chan struct{}) {
	streams := make([]<-chan Allocation, 0, len(producers))
	errs := make(chan error, len(producers))
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(producers))
	for _, producer := range producers {
		out := make(chan Allocation)
		streams = append(streams, out)
		current := producer
		go func() {
			defer wg.Done()
			current.Stream(out, errs)
		}()
	}
	go func() {
		wg.Wait()
		close(errs)
		close(done)
	}()
	return streams, errs, done
}
