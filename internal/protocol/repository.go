package protocol

import (
	"errors"
	"fmt"
	"sync"
)

var ErrDeviceMissing = errors.New("protocol device missing")

type Repository struct {
	mu     sync.Mutex
	frames map[uint32][]Frame
}

func NewRepository() *Repository { return &Repository{frames: make(map[uint32][]Frame)} }

func (r *Repository) Store(frame Frame) error {
	if frame.DeviceID == 0 {
		return fmt.Errorf("store device: %w", ErrDeviceMissing)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frames[frame.DeviceID] = append(r.frames[frame.DeviceID], frame)
	return nil
}

func (r *Repository) Count(device uint32) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.frames[device])
}
