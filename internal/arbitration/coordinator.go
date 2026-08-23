package arbitration

import (
	"context"
	"errors"
	"sync"
)

var ErrCoordinatorClosed = errors.New("arbitration coordinator closed")

type finishRequest struct {
	id       string
	state    State
	response chan error
}

type Coordinator struct {
	registry *Registry
	ctx      context.Context
	cancel   context.CancelFunc
	jobs     chan finishRequest
	mu       sync.RWMutex
	closed   bool
	wg       sync.WaitGroup
	once     sync.Once
}

func NewCoordinator(registry *Registry, queue int) *Coordinator {
	if queue < 1 {
		queue = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Coordinator{registry: registry, ctx: ctx, cancel: cancel, jobs: make(chan finishRequest, queue)}
	c.wg.Add(1)
	go c.run()
	return c
}

func (c *Coordinator) Submit(ctx context.Context, id string, state State) error {
	req := finishRequest{id: id, state: state, response: make(chan error, 1)}
	c.mu.RLock()
	c.jobs <- req
	c.mu.RUnlock()
	select {
	case err := <-req.response:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-c.ctx.Done():
		return ErrCoordinatorClosed
	}
}

func (c *Coordinator) Close() {
	c.once.Do(func() {
		c.cancel()
		c.mu.Lock()
		c.closed = true
		c.mu.Unlock()
		c.wg.Wait()
		for {
			select {
			case req := <-c.jobs:
				req.response <- ErrCoordinatorClosed
			default:
				return
			}
		}
	})
}

func (c *Coordinator) run() {
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		case req := <-c.jobs:
			req.response <- c.registry.Finish(req.id, req.state)
		}
	}
}
