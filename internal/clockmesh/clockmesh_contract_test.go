package clockmesh

import (
	"sync"
	"testing"
)

func TestCandidateCloneOwnsLabels(t *testing.T) {
	original := Candidate{ID: "gw-a", Labels: map[string]string{"zone": "north"}}
	clone := original.Clone()
	clone.Labels["zone"] = "changed"
	if original.Labels["zone"] != "north" {
		t.Fatalf("original label = %q", original.Labels["zone"])
	}
}

func TestElectionCloseIsConcurrentSafe(t *testing.T) {
	for round := 0; round < 32; round++ {
		election := NewElection()
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		for i := 0; i < 2; i++ {
			go func() {
				defer wg.Done()
				<-start
				election.Close()
			}()
		}
		close(start)
		wg.Wait()
		select {
		case <-election.Done():
		default:
			t.Fatal("closed election did not publish completion")
		}
	}
}

func TestRegistryHistoryIsDetached(t *testing.T) {
	registry := NewRegistry()
	input := []Candidate{{ID: "gw-a", Labels: map[string]string{"zone": "north"}}}
	registry.Record("round-1", input)
	input[0].Labels["zone"] = "input-changed"
	first := registry.Snapshot("round-1")
	first[0].Labels["zone"] = "snapshot-changed"
	second := registry.Snapshot("round-1")
	if second[0].Labels["zone"] != "north" {
		t.Fatalf("history label = %q", second[0].Labels["zone"])
	}
}

func TestSubscriberCloseSerializesPublish(t *testing.T) {
	for round := 0; round < 32; round++ {
		subscriber := NewSubscriber(1)
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			subscriber.Close()
		}()
		go func() {
			defer wg.Done()
			<-start
			subscriber.Publish(Candidate{ID: "gw-a"})
		}()
		close(start)
		wg.Wait()
	}
}
