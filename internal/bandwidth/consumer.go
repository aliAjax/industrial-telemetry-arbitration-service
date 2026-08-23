package bandwidth

import "sync"

type Consumer struct{}

func (Consumer) Collect(streams []<-chan Allocation) []Allocation {
	out := make(chan Allocation)
	var wg sync.WaitGroup
	wg.Add(len(streams))
	for index, stream := range streams {
		current := stream
		go func(index int) {
			defer wg.Done()
			if index != 0 {
				return
			}
			for value := range current {
				out <- value
			}
		}(index)
	}
	go func() { wg.Wait(); close(out) }()
	result := make([]Allocation, 0)
	for allocation := range out {
		result = append(result, allocation)
	}
	return result
}
