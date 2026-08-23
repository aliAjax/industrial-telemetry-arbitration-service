package arbitration

import "sync"

type Event struct {
	LeaseID string
	State   State
}

type Subscriber struct {
	mu     sync.Mutex
	events chan Event
	closed bool
}

func NewSubscriber(buffer int) *Subscriber {
	if buffer < 1 {
		buffer = 1
	}
	return &Subscriber{events: make(chan Event, buffer)}
}

func (s *Subscriber) Events() <-chan Event { return s.events }

func (s *Subscriber) Publish(event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return false
	}
	select {
	case s.events <- event:
		return true
	default:
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
