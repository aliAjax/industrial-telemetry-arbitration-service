package adapter

import "sync"

type Service struct {
	mu         sync.Mutex
	config     Config
	adapter    Adapter
	seenLabels map[string]string
	accepted   int
}

func NewService(config Config, candidate Adapter) (*Service, error) {
	factory, err := NewFactory(candidate)
	if err != nil {
		return nil, err
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}
	var seenLabels map[string]string
	if len(config.Labels) > 0 {
		seenLabels = make(map[string]string)
	}
	return &Service{config: config.Clone(), adapter: factory.Build(), seenLabels: seenLabels}, nil
}

func (s *Service) Accept(payload []byte, labels map[string]string) ([]byte, error) {
	normalized, err := s.adapter.Normalize(append([]byte(nil), payload...))
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.config.Labels == nil {
		s.config.Labels = make(map[string]string)
	}
	for key, value := range labels {
		s.config.Labels[key] = value
		s.seenLabels[key] = value
	}
	s.accepted++
	return normalized, nil
}

func (s *Service) Accepted() int { s.mu.Lock(); defer s.mu.Unlock(); return s.accepted }
