package transaction

import (
	"database/sql"

	"otas/models"
)

type transactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{DB: db}
}

// for paths 1 and 2; no DB transaction needed
func (r *transactionRepository) CreateTransaction(t *models.Transaction) (*models.Transaction, error) {
    query := `
        INSERT INTO transactions (user_id, amount, deduction, allocated_to, type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, user_id, amount, deduction, allocated_to, type, created_at
    `
    result := &models.Transaction{}
    err := r.DB.QueryRow(query,
        t.UserID, t.Amount, t.Deduction, t.AllocatedTo, t.Type,
    ).Scan(
        &result.ID, &result.UserID, &result.Amount,
        &result.Deduction, &result.AllocatedTo, &result.Type, &result.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return result, nil
}

//for path 3, Tx needed
func (r *transactionRepository) CreateTransactionTx(tx *sql.Tx, t *models.Transaction) (*models.Transaction, error) {
	query := `
        INSERT INTO transactions (user_id, amount, deduction, allocated_to, type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, user_id, amount, deduction, allocated_to, type, created_at
    `
	result := &models.Transaction{}
	err := tx.QueryRow(query,
		t.UserID,
		t.Amount,
		t.Deduction,
		t.AllocatedTo,
		t.Type,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Amount,
		&result.Deduction,
		&result.AllocatedTo,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// count how many eligible transactions user has made today
func (r *transactionRepository) CountTodayTransactions(userID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM transactions
		WHERE user_id = $1
		AND created_at >= CURRENT_DATE
	`
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// get total contributed by a user,  for group contribution tracking
func (r *transactionRepository) GetTotalContributed(userID int, accountType string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(deduction), 0)
		FROM transactions
		WHERE user_id = $1
		AND allocated_to = $2
	`
	var total float64
	err := r.DB.QueryRow(query, userID, accountType).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
