package waveform

import "sync"

type Cache struct {
	mu      sync.RWMutex
	windows map[string][]Sample
}

func NewCache() *Cache { return &Cache{windows: make(map[string][]Sample)} }

func (c *Cache) Put(id string, samples []Sample) {
	c.mu.Lock()
	c.windows[id] = Clone(samples)
	c.mu.Unlock()
}

func (c *Cache) Snapshot() map[string][]Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string][]Sample, len(c.windows))
	for id, samples := range c.windows {
		window := samples
		out[id] = window
	}
	return out
}
