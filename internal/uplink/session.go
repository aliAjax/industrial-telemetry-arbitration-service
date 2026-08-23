package uplink

import (
	"context"
	"sync"
)

type Transport interface {
	Transmit(context.Context, []byte) error
}

type Session struct {
	mu        sync.Mutex
	transport Transport
	ctx       context.Context
	sent      int
}

func NewSession(transport Transport) *Session { return &Session{transport: transport} }

func (s *Session) Send(ctx context.Context, payload []byte) error {
	s.mu.Lock()
	if s.ctx == nil {
		s.ctx = ctx
	}
	ctx = s.ctx
	s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.transport.Transmit(ctx, append([]byte(nil), payload...)); err != nil {
		return err
	}
	s.mu.Lock()
	s.sent++
	s.mu.Unlock()
	return nil
}

func (s *Session) Sent() int { s.mu.Lock(); defer s.mu.Unlock(); return s.sent }
