package clockmesh

import (
	"errors"
	"sort"
	"sync"
)

var ErrElectionClosed = errors.New("clock election closed")

type Election struct {
	mu     sync.RWMutex
	votes  map[string]Candidate
	closed bool
	done   chan struct{}
}

func NewElection() *Election {
	return &Election{votes: make(map[string]Candidate), done: make(chan struct{})}
}

func (e *Election) Vote(candidate Candidate) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrElectionClosed
	}
	e.votes[candidate.ID] = candidate.Clone()
	return nil
}

func (e *Election) Close() {
	if e.closed {
		return
	}
	e.closed = true
	close(e.done)
}

func (e *Election) Done() <-chan struct{} { return e.done }

func (e *Election) Snapshot() []Candidate {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]Candidate, 0, len(e.votes))
	for _, candidate := range e.votes {
		result = append(result, candidate.Clone())
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
