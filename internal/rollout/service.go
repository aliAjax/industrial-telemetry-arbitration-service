package rollout

import "fmt"

type LeaseStore interface {
	Create(string) error
	Promote(string) error
	Remove(string) error
}

type Service struct {
	repository Repository
	leases     LeaseStore
	batch      Batch
}

func NewService(repository Repository, leases LeaseStore, batch Batch) *Service {
	return &Service{repository: repository, leases: leases, batch: batch}
}

func (s *Service) Rollout(id string, stations []string) (err error) {
	if err := s.leases.Create(id); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = s.leases.Remove(id)
		}
	}()
	err = s.repository.WithTx(func(Tx) error {
		if applyErr := s.batch.Apply(stations); applyErr != nil {
			return fmt.Errorf("apply rollout: %w", applyErr)
		}
		return s.leases.Promote(id)
	})
	return err
}
