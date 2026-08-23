package rollout

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidRolloutPolicy = errors.New("invalid rollout policy")

type RolloutPolicy struct {
	ID        string
	Station   string
	Batch     int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p RolloutPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Station) == "" {
		return ErrInvalidRolloutPolicy
	}
	if p.Batch < 0 {
		return ErrInvalidRolloutPolicy
	}
	return nil
}

func (p RolloutPolicy) Copy() RolloutPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type RolloutPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]RolloutPolicy
	revision uint64
}

func NewRolloutPolicyCatalog() *RolloutPolicyCatalog {
	return &RolloutPolicyCatalog{byID: make(map[string]RolloutPolicy)}
}

func (c *RolloutPolicyCatalog) Put(policy RolloutPolicy) error {
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

func (c *RolloutPolicyCatalog) Get(id string) (RolloutPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *RolloutPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *RolloutPolicyCatalog) List(enabledOnly bool) []RolloutPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]RolloutPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *RolloutPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *RolloutPolicyCatalog) Replace(policies []RolloutPolicy) error {
	next := make(map[string]RolloutPolicy, len(policies))
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

func (c *RolloutPolicyCatalog) CountByLabel(key, value string) int {
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
