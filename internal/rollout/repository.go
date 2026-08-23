package rollout

import (
	"errors"
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
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
	}()
	if err = operation(tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit rollout: %w", err)
	}
	committed = true
	return nil
}
