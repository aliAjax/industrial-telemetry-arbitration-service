package rollout

import (
	"errors"
	"fmt"
)

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
	promoted := false
	defer func() {
		if promoted {
			return
		}
		if rmErr := s.leases.Remove(id); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("remove lease %s: %w", id, rmErr))
		}
	}()

	err = s.repository.WithTx(func(Tx) error {
		if applyErr := s.batch.Apply(stations); applyErr != nil {
			return fmt.Errorf("apply rollout: %w", applyErr)
		}
		if promoteErr := s.leases.Promote(id); promoteErr != nil {
			return fmt.Errorf("promote lease %s: %w", id, promoteErr)
		}
		return nil
	})
	if err == nil {
		promoted = true
	}
	return err
}
