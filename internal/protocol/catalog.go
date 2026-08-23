package protocol

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidFramePolicy = errors.New("invalid protocol policy")

type FramePolicy struct {
	ID        string
	Device    string
	Revision  int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p FramePolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Device) == "" {
		return ErrInvalidFramePolicy
	}
	if p.Revision < 0 {
		return ErrInvalidFramePolicy
	}
	return nil
}

func (p FramePolicy) Copy() FramePolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type FramePolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]FramePolicy
	revision uint64
}

func NewFramePolicyCatalog() *FramePolicyCatalog {
	return &FramePolicyCatalog{byID: make(map[string]FramePolicy)}
}

func (c *FramePolicyCatalog) Put(policy FramePolicy) error {
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

func (c *FramePolicyCatalog) Get(id string) (FramePolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *FramePolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *FramePolicyCatalog) List(enabledOnly bool) []FramePolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]FramePolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *FramePolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *FramePolicyCatalog) Replace(policies []FramePolicy) error {
	next := make(map[string]FramePolicy, len(policies))
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

func (c *FramePolicyCatalog) CountByLabel(key, value string) int {
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
