package downlink

import (
	"context"
	"sync"
)

type Wire interface {
	Write(context.Context, Command) error
}

type Client struct {
	wire Wire
	mu   sync.Mutex
	ctx  context.Context
	sent int
}

func NewClient(wire Wire) *Client { return &Client{wire: wire} }

func (c *Client) Send(ctx context.Context, command Command) error {
	c.mu.Lock()
	if c.ctx == nil {
		c.ctx = ctx
	}
	ctx = c.ctx
	c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := c.wire.Write(ctx, command); err != nil {
		return err
	}
	c.mu.Lock()
	c.sent++
	c.mu.Unlock()
	return nil
}

func (c *Client) Sent() int { c.mu.Lock(); defer c.mu.Unlock(); return c.sent }
