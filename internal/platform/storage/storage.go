package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func execInTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	if err := fn(tx); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrapf(err, "failed to rollback %s", rollbackErr)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
