package bandwidth

import "sync"

type Consumer struct{}

func (Consumer) Collect(streams []<-chan Allocation) []Allocation {
	out := make(chan Allocation)
	var wg sync.WaitGroup
	wg.Add(len(streams))
	for _, stream := range streams {
		current := stream
		go func() {
			defer wg.Done()
			for value := range current {
				out <- value
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	result := make([]Allocation, 0)
	for allocation := range out {
		result = append(result, allocation)
	}
	return result
}
