package ackgroup

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPublisherCoordinatesShutdown(t *testing.T) {
	for round := 0; round < 32; round++ {
		publisher := NewPublisher(1)
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			publisher.Close()
		}()
		go func() {
			defer wg.Done()
			<-start
			publisher.Publish(Ack{GatewayID: "gw-a"})
		}()
		close(start)
		wg.Wait()
	}
}

func TestFanoutJoinsBeforeClosing(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	processor := func(_ context.Context, ack Ack) (Ack, error) {
		close(started)
		<-release
		return ack, nil
	}
	results, errorsOut := Fanout(context.Background(), []Ack{{GatewayID: "gw-a"}}, processor)
	<-started
	close(release)
	var collected []Ack
	for ack := range results {
		collected = append(collected, ack)
	}
	for err := range errorsOut {
		if err != nil {
			t.Fatalf("fanout error: %v", err)
		}
	}
	if len(collected) != 1 {
		t.Fatalf("collected %d acknowledgements", len(collected))
	}
}

func TestCollectorDrainsAllAcks(t *testing.T) {
	results := make(chan Ack, 3)
	results <- Ack{GatewayID: "gw-a"}
	results <- Ack{GatewayID: "gw-b"}
	results <- Ack{GatewayID: "gw-c"}
	close(results)
	if collected := Collect(results); len(collected) != 3 {
		t.Fatalf("collected %d acknowledgements", len(collected))
	}
}

func TestStoreSnapshotOwnsMetadata(t *testing.T) {
	store := &Store{}
	input := Ack{GatewayID: "gw-a", Metadata: map[string]string{"zone": "north"}}
	store.Append(input)
	input.Metadata["zone"] = "input-changed"
	first := store.Snapshot()
	first[0].Metadata["zone"] = "snapshot-changed"
	second := store.Snapshot()
	if second[0].Metadata["zone"] != "north" {
		t.Fatalf("stored metadata = %q", second[0].Metadata["zone"])
	}
	time.Sleep(time.Millisecond)
}
