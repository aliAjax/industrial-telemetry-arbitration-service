package clockmesh

import "sync"

type Registry struct {
	mu      sync.RWMutex
	history map[string][]Candidate
}

func NewRegistry() *Registry {
	return &Registry{history: make(map[string][]Candidate)}
}

func (r *Registry) Record(key string, candidates []Candidate) {
	r.mu.Lock()
	r.history[key] = candidates
	r.mu.Unlock()
}

func (r *Registry) Snapshot(key string) []Candidate {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.history[key]
}

func cloneCandidates(candidates []Candidate) []Candidate {
	result := make([]Candidate, len(candidates))
	for i, candidate := range candidates {
		result[i] = candidate.Clone()
	}
	return result
}
