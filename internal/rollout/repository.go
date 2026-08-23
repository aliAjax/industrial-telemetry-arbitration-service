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
	committed := false
	defer func() {
		if committed {
			return
		}
		if rbErr := tx.Rollback(); rbErr != nil && err == nil {
			err = fmt.Errorf("rollback rollout: %w", rbErr)
		}
	}()
	if err = operation(tx); err != nil {
		return err
	}
	if cerr := tx.Commit(); cerr != nil {
		err = fmt.Errorf("commit rollout: %w", cerr)
		return err
	}
	committed = true
	return nil
}
