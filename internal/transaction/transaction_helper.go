package transaction

import "database/sql"

func (r *transactionRepository) WithTx(fn func(*sql.Tx) error) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
