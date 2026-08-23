package downlink

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidDispatchPolicy = errors.New("invalid downlink policy")

type DispatchPolicy struct {
	ID        string
	Device    string
	Priority  int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p DispatchPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Device) == "" {
		return ErrInvalidDispatchPolicy
	}
	if p.Priority < 0 {
		return ErrInvalidDispatchPolicy
	}
	return nil
}

func (p DispatchPolicy) Copy() DispatchPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type DispatchPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]DispatchPolicy
	revision uint64
}

func NewDispatchPolicyCatalog() *DispatchPolicyCatalog {
	return &DispatchPolicyCatalog{byID: make(map[string]DispatchPolicy)}
}

func (c *DispatchPolicyCatalog) Put(policy DispatchPolicy) error {
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

func (c *DispatchPolicyCatalog) Get(id string) (DispatchPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *DispatchPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *DispatchPolicyCatalog) List(enabledOnly bool) []DispatchPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]DispatchPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *DispatchPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *DispatchPolicyCatalog) Replace(policies []DispatchPolicy) error {
	next := make(map[string]DispatchPolicy, len(policies))
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

func (c *DispatchPolicyCatalog) CountByLabel(key, value string) int {
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
