package arbitration

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidLeasePolicy = errors.New("invalid arbitration policy")

type LeasePolicy struct {
	ID        string
	Gateway   string
	Priority  int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p LeasePolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Gateway) == "" {
		return ErrInvalidLeasePolicy
	}
	if p.Priority < 0 {
		return ErrInvalidLeasePolicy
	}
	return nil
}

func (p LeasePolicy) Copy() LeasePolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type LeasePolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]LeasePolicy
	revision uint64
}

func NewLeasePolicyCatalog() *LeasePolicyCatalog {
	return &LeasePolicyCatalog{byID: make(map[string]LeasePolicy)}
}

func (c *LeasePolicyCatalog) Put(policy LeasePolicy) error {
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

func (c *LeasePolicyCatalog) Get(id string) (LeasePolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *LeasePolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *LeasePolicyCatalog) List(enabledOnly bool) []LeasePolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]LeasePolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *LeasePolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *LeasePolicyCatalog) Replace(policies []LeasePolicy) error {
	next := make(map[string]LeasePolicy, len(policies))
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

func (c *LeasePolicyCatalog) CountByLabel(key, value string) int {
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
