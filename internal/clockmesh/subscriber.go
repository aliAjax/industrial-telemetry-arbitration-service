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
	if s.closed {
		return false
	}
	select {
	case s.events <- candidate.Clone():
		return true
	default:
		return false
	}
}

func (s *Subscriber) Close() {
	if s.closed {
		return
	}
	s.closed = true
	close(s.events)
}
