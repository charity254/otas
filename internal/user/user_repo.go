// internal/user/user_repo.go
package user

import (
	"database/sql"

	"otas/models"
)

type userRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *userRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, phone, password, saving_type, daily_limit)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, name, email, phone, saving_type, daily_limit, created_at, updated_at
    `
	result := &models.User{}
	err := r.DB.QueryRow(query,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.SavingType,
		user.DailyLimit,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.SavingType,
		&result.DailyLimit,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, name, email, phone, password, saving_type, daily_limit, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.SavingType,
		&user.DailyLimit,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
