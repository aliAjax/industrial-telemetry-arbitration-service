package ackgroup

import "sync"

type Store struct {
	mu   sync.RWMutex
	acks []Ack
}

func (s *Store) Append(ack Ack) {
	s.mu.Lock()
	s.acks = append(s.acks, ack.Clone())
	s.mu.Unlock()
}

func (s *Store) Snapshot() []Ack {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Ack, len(s.acks))
	for i, ack := range s.acks {
		result[i] = ack.Clone()
	}
	return result
}
