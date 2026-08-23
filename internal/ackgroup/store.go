package ackgroup

import "sync"

type Store struct {
	mu   sync.RWMutex
	acks []Ack
}

func (s *Store) Append(ack Ack) {
	clone := ack.Clone()
	s.mu.Lock()
	s.acks = append(s.acks, clone)
	s.mu.Unlock()
}

func (s *Store) Snapshot() []Ack {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cloned := make([]Ack, len(s.acks))
	for i, ack := range s.acks {
		cloned[i] = ack.Clone()
	}
	return cloned
}
