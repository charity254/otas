// internal/account/account_repo.go
package account

import (
	"database/sql"
	"otas/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) CreateAccount(acct *models.Account) (*models.Account, error) {
	query := `
        INSERT INTO accounts (user_id, type, balance)
        VALUES ($1, $2, 0.00)
        RETURNING id, user_id, type, balance, created_at, updated_at
    `
	account := &models.Account{}
	err := r.DB.QueryRow(query,
		acct.UserID,
		acct.Type,
		acct.Balance,
		acct.CreatedAt,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Type,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) GetAccountsByUserID(userID int) ([]models.Account, error) {
	query := `
        SELECT id, user_id, type, balance, created_at, updated_at
        FROM accounts
        WHERE user_id = $1
    `
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Type,
			&a.Balance,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(userID int, accountType models.AccountType, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = NOW()
		WHERE user_id = $2
		AND type = $3
	`
	_, err := r.DB.Exec(query, amount, userID, accountType)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) GetAccountBalance(userID int, accountType models.AccountType) (float64, error) {
	query := `
		SELECT balance
		FROM accounts
		WHERE user_id = $1
		AND type = $2
	`
	var balance float64
	err := r.DB.QueryRow(query, userID, accountType).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
