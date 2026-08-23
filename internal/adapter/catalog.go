package adapter

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidAdapterPolicy = errors.New("invalid adapter policy")

type AdapterPolicy struct {
	ID        string
	Vendor    string
	Weight    int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p AdapterPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Vendor) == "" {
		return ErrInvalidAdapterPolicy
	}
	if p.Weight < 0 {
		return ErrInvalidAdapterPolicy
	}
	return nil
}

func (p AdapterPolicy) Copy() AdapterPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type AdapterPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]AdapterPolicy
	revision uint64
}

func NewAdapterPolicyCatalog() *AdapterPolicyCatalog {
	return &AdapterPolicyCatalog{byID: make(map[string]AdapterPolicy)}
}

func (c *AdapterPolicyCatalog) Put(policy AdapterPolicy) error {
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

func (c *AdapterPolicyCatalog) Get(id string) (AdapterPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *AdapterPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *AdapterPolicyCatalog) List(enabledOnly bool) []AdapterPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]AdapterPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *AdapterPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *AdapterPolicyCatalog) Replace(policies []AdapterPolicy) error {
	next := make(map[string]AdapterPolicy, len(policies))
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

func (c *AdapterPolicyCatalog) CountByLabel(key, value string) int {
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
