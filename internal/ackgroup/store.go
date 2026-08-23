package ackgroup

import "sync"

type Store struct {
	mu   sync.RWMutex
	acks []Ack
}

func (s *Store) Append(ack Ack) {
	s.mu.Lock()
	s.acks = append(s.acks, ack)
	s.mu.Unlock()
}

func (s *Store) Snapshot() []Ack {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.acks
}
