package rollout

import (
	"fmt"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository struct{ Begin func() (Tx, error) }

func (r Repository) WithTx(operation func(Tx) error) (err error) {
	tx, err := r.Begin()
	if err != nil {
		return fmt.Errorf("begin rollout: %w", err)
	}
	defer func() { err = tx.Commit() }()
	if err = operation(tx); err != nil {
		return err
	}
	return nil
}
