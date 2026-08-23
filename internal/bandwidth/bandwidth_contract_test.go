package bandwidth

import (
	"sync"
	"testing"
	"time"
)

func TestProducerClosesOutputOnRejectedLink(t *testing.T) {
	out := make(chan Allocation)
	errs := make(chan error, 1)
	start := make(chan struct{})
	done := make(chan struct{})
	go func() { <-start; (Producer{Reject: true}).Stream(out, errs); close(done) }()
	close(start)
	<-done
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("unexpected allocation")
		}
	default:
		t.Fatal("rejected producer left output open")
	}
}

func TestCoordinatorWaitsForAllProducers(t *testing.T) {
	streams, _, done := (Coordinator{}).Rebalance([]Producer{{LinkID: "a", Units: []int{1}}})
	select {
	case <-done:
		t.Fatal("coordinator completed before producer output was consumed")
	case <-time.After(20 * time.Millisecond):
	}
	go func() {
		for range streams[0] {
		}
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("coordinator did not finish")
	}
}

func TestConsumerStopsAfterAllStreamsClose(t *testing.T) {
	first := make(chan Allocation, 1)
	second := make(chan Allocation, 1)
	first <- Allocation{LinkID: "a"}
	second <- Allocation{LinkID: "b"}
	close(first)
	close(second)
	start := make(chan struct{})
	var values []Allocation
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); <-start; values = (Consumer{}).Collect([]<-chan Allocation{first, second}) }()
	close(start)
	wg.Wait()
	if len(values) != 2 {
		t.Fatalf("allocations = %#v", values)
	}
}

func TestResultErrorDoesNotBlockWorkers(t *testing.T) {
	result := NewResult(2)
	start := make(chan struct{})
	done := make(chan bool, 1)
	go func() { <-start; done <- result.PublishError(ErrLinkRejected) }()
	close(start)
	select {
	case ok := <-done:
		if !ok {
			t.Fatal("buffered error was dropped")
		}
	case <-time.After(time.Second):
		t.Fatal("error publisher blocked")
	}
}
