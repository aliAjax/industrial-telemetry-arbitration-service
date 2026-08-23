package clockmesh

import "sync"

type Subscriber struct {
	mu     sync.Mutex
	events chan Candidate
	closed bool
}

func NewSubscriber(capacity int) *Subscriber {
	return &Subscriber{events: make(chan Candidate, capacity)}
}

func (s *Subscriber) Events() <-chan Candidate { return s.events }

func (s *Subscriber) Publish(candidate Candidate) bool {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return false
	}
	select {
	case s.events <- candidate.Clone():
		s.mu.Unlock()
		return true
	default:
		s.mu.Unlock()
		return false
	}
}

func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	close(s.events)
}
