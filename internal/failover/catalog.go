package failover

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidRecoveryPolicy = errors.New("invalid failover policy")

type RecoveryPolicy struct {
	ID        string
	Link      string
	Threshold int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p RecoveryPolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Link) == "" {
		return ErrInvalidRecoveryPolicy
	}
	if p.Threshold < 0 {
		return ErrInvalidRecoveryPolicy
	}
	return nil
}

func (p RecoveryPolicy) Copy() RecoveryPolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type RecoveryPolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]RecoveryPolicy
	revision uint64
}

func NewRecoveryPolicyCatalog() *RecoveryPolicyCatalog {
	return &RecoveryPolicyCatalog{byID: make(map[string]RecoveryPolicy)}
}

func (c *RecoveryPolicyCatalog) Put(policy RecoveryPolicy) error {
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

func (c *RecoveryPolicyCatalog) Get(id string) (RecoveryPolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *RecoveryPolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *RecoveryPolicyCatalog) List(enabledOnly bool) []RecoveryPolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]RecoveryPolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *RecoveryPolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *RecoveryPolicyCatalog) Replace(policies []RecoveryPolicy) error {
	next := make(map[string]RecoveryPolicy, len(policies))
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

func (c *RecoveryPolicyCatalog) CountByLabel(key, value string) int {
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
