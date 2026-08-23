package bandwidth

import "sync"

type Result struct {
	mu     sync.Mutex
	Errors []error
	notify chan error
}

func NewResult(capacity int) *Result {
	if capacity < 1 {
		capacity = 1
	}
	return &Result{notify: make(chan error, capacity)}
}

func (r *Result) PublishError(err error) bool {
	select {
	case r.notify <- err:
		r.mu.Lock()
		r.Errors = append(r.Errors, err)
		r.mu.Unlock()
		return true
	default:
		return false
	}
}

func (r *Result) Notifications() <-chan error { return r.notify }
func (r *Result) Count() int                  { r.mu.Lock(); defer r.mu.Unlock(); return len(r.Errors) }
