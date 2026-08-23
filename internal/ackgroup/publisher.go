package ackgroup

import "sync"

type Publisher struct {
	mu     sync.Mutex
	events chan Ack
	closed bool
}

func NewPublisher(capacity int) *Publisher {
	return &Publisher{events: make(chan Ack, capacity)}
}

func (p *Publisher) Events() <-chan Ack { return p.events }

func (p *Publisher) Publish(ack Ack) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return false
	}
	select {
	case p.events <- ack.Clone():
		return true
	default:
		return false
	}
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	p.closed = true
	close(p.events)
}
