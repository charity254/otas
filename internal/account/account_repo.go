// internal/account/account_repo.go
package account

import (
    "database/sql"
    "otas/models"
)

type Repository struct {
    DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{DB: db}
}

func (r *Repository) CreateAccount(userID int, accountType models.AccountType) (*models.Account, error) {
    query := `
        INSERT INTO accounts (user_id, type, balance)
        VALUES ($1, $2, 0.00)
        RETURNING id, user_id, type, balance, created_at, updated_at
    `
    account := &models.Account{}
    err := r.DB.QueryRow(query, userID, accountType).Scan(
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

func (r *Repository) GetAccountsByUserID(userID int) ([]models.Account, error) {
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

