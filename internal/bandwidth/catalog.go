package bandwidth

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidCapacityPolicy = errors.New("invalid bandwidth policy")

type CapacityPolicy struct {
	ID        string
	Link      string
	Units     int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p CapacityPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Link) == "" {
		return ErrInvalidCapacityPolicy
	}
	if p.Units < 0 {
		return ErrInvalidCapacityPolicy
	}
	return nil
}

func (p CapacityPolicy) Copy() CapacityPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type CapacityPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]CapacityPolicy
	revision uint64
}

func NewCapacityPolicyCatalog() *CapacityPolicyCatalog {
	return &CapacityPolicyCatalog{byID: make(map[string]CapacityPolicy)}
}

func (c *CapacityPolicyCatalog) Put(policy CapacityPolicy) error {
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

func (c *CapacityPolicyCatalog) Get(id string) (CapacityPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *CapacityPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *CapacityPolicyCatalog) List(enabledOnly bool) []CapacityPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]CapacityPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *CapacityPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *CapacityPolicyCatalog) Replace(policies []CapacityPolicy) error {
	next := make(map[string]CapacityPolicy, len(policies))
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

func (c *CapacityPolicyCatalog) CountByLabel(key, value string) int {
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
