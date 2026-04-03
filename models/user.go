package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Password   string    `json:"-"`
	SavingType string    `json:"saving_type"`
	DailyLimit int       `json:"daily_limit"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
