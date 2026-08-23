package waveform

import "sync"

type Archive struct {
	mu      sync.RWMutex
	windows map[string][]Sample
}

func NewArchive() *Archive { return &Archive{windows: make(map[string][]Sample)} }

func (a *Archive) Store(id string, samples []Sample) {
	a.mu.Lock()
	a.windows[id] = Clone(samples)
	a.mu.Unlock()
}

func (a *Archive) Load(id string) []Sample {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return Clone(a.windows[id])
}
