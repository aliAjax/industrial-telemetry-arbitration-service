package arbitration

import (
	"errors"
	"sort"
	"sync"
)

var ErrLeaseMissing = errors.New("arbitration lease missing")

type Registry struct {
	mu     sync.RWMutex
	leases map[string]*Lease
}

func NewRegistry() *Registry { return &Registry{leases: make(map[string]*Lease)} }

func (r *Registry) Create(id string, metadata map[string]string) *Lease {
	r.mu.Lock()
	defer r.mu.Unlock()
	lease := NewLease(id, metadata)
	r.leases[id] = lease
	return lease
}

func (r *Registry) Finish(id string, state State) error {
	r.mu.RLock()
	lease := r.leases[id]
	r.mu.RUnlock()
	if lease == nil {
		return ErrLeaseMissing
	}
	return lease.Finish(state)
}

func (r *Registry) Snapshot() []LeaseSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]LeaseSnapshot, 0, len(r.leases))
	for _, lease := range r.leases {
		out = append(out, LeaseSnapshot{ID: lease.id, State: lease.state, Metadata: lease.metadata})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
