package uplink

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrInvalidRoutePolicy = errors.New("invalid uplink policy")

type RoutePolicy struct {
	ID        string
	Gateway   string
	Window    int
	Enabled   bool
	Labels    map[string]string
	UpdatedAt time.Time
}

func (p RoutePolicy) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Gateway) == "" {
		return ErrInvalidRoutePolicy
	}
	if p.Window < 0 {
		return ErrInvalidRoutePolicy
	}
	return nil
}

func (p RoutePolicy) Copy() RoutePolicy {
	out := p
	out.Labels = make(map[string]string, len(p.Labels))
	for key, value := range p.Labels {
		out.Labels[key] = value
	}
	return out
}

type RoutePolicyCatalog struct {
	mu       sync.RWMutex
	byID     map[string]RoutePolicy
	revision uint64
}

func NewRoutePolicyCatalog() *RoutePolicyCatalog {
	return &RoutePolicyCatalog{byID: make(map[string]RoutePolicy)}
}

func (c *RoutePolicyCatalog) Put(policy RoutePolicy) error {
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

func (c *RoutePolicyCatalog) Get(id string) (RoutePolicy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.byID[id]
	return policy.Copy(), ok
}

func (c *RoutePolicyCatalog) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.byID[id]; !ok {
		return false
	}
	delete(c.byID, id)
	c.revision++
	return true
}

func (c *RoutePolicyCatalog) List(enabledOnly bool) []RoutePolicy {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]RoutePolicy, 0, len(c.byID))
	for _, policy := range c.byID {
		if enabledOnly && !policy.Enabled {
			continue
		}
		out = append(out, policy.Copy())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (c *RoutePolicyCatalog) Revision() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.revision
}

func (c *RoutePolicyCatalog) Replace(policies []RoutePolicy) error {
	next := make(map[string]RoutePolicy, len(policies))
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

func (c *RoutePolicyCatalog) CountByLabel(key, value string) int {
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
