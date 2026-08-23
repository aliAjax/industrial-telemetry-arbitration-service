package arbitration

import (
	"context"
	"errors"
	"sync"
)

var ErrLeaseFinished = errors.New("arbitration lease already finished")

type State string

const (
	StateActive State = "active"
	StateWon    State = "won"
	StateLost   State = "lost"
)

type LeaseSnapshot struct {
	ID       string
	State    State
	Metadata map[string]string
}

type Lease struct {
	mu       sync.RWMutex
	id       string
	state    State
	metadata map[string]string
	done     chan struct{}
}

func NewLease(id string, metadata map[string]string) *Lease {
	return &Lease{id: id, state: StateActive, metadata: cloneMap(metadata), done: make(chan struct{})}
}

func (l *Lease) Finish(next State) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.state != StateActive {
		return ErrLeaseFinished
	}
	l.state = next
	close(l.done)
	return nil
}

func (l *Lease) Wait(ctx context.Context) (State, error) {
	select {
	case <-ctx.Done():
		return StateActive, ctx.Err()
	case <-l.done:
		return l.Snapshot().State, nil
	}
}

func (l *Lease) Snapshot() LeaseSnapshot {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return LeaseSnapshot{ID: l.id, State: l.state, Metadata: cloneMap(l.metadata)}
}

func cloneMap(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}
