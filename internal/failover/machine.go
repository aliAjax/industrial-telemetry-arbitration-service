package failover

import (
	"errors"
	"sync"
)

var ErrIllegalTransition = errors.New("illegal failover transition")

type Machine struct {
	mu      sync.RWMutex
	state   State
	history []State
}

func NewMachine(initial State) *Machine { return &Machine{state: initial, history: []State{initial}} }

func (m *Machine) Advance(next State) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.state.CanTransition(next) {
		return ErrIllegalTransition
	}
	m.state = next
	m.history = append(m.history, next)
	return nil
}

func (m *Machine) State() State { m.mu.RLock(); defer m.mu.RUnlock(); return m.state }
func (m *Machine) History() []State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]State(nil), m.history...)
}
