package waveform

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidWindowPolicy = errors.New("invalid waveform policy")

type WindowPolicy struct {
	ID        string
	Sensor    string
	Samples   int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p WindowPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Sensor) == "" {
		return ErrInvalidWindowPolicy
	}
	if p.Samples < 0 {
		return ErrInvalidWindowPolicy
	}
	return nil
}

func (p WindowPolicy) Copy() WindowPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type WindowPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]WindowPolicy
	revision uint64
}

func NewWindowPolicyCatalog() *WindowPolicyCatalog {
	return &WindowPolicyCatalog{byID: make(map[string]WindowPolicy)}
}

func (c *WindowPolicyCatalog) Put(policy WindowPolicy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	policy.UpdatedAt = policy.UpdatedAt.UTC()
	c.byID[policy.ID] = policy.Copy()
	c.revision++
	return nil
}

func (c *WindowPolicyCatalog) Get(id string) (WindowPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *WindowPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *WindowPolicyCatalog) List(enabledOnly bool) []WindowPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]WindowPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *WindowPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *WindowPolicyCatalog) Replace(policies []WindowPolicy) error {
	next := make(map[string]WindowPolicy, len(policies))
	for _, policy := range policies {
		if err := policy.Validate(); err != nil {
			return err
		}
		next[policy.ID] = policy.Copy()
	}
	c.mu.Lock()
	c.byID = next
	c.revision++
	c.mu.Unlock()
	return nil
}

func (c *WindowPolicyCatalog) CountByLabel(key, value string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, policy := range c.byID {
		if policy.Labels[key] == value {
			count++
		}
	}
	return count
}
