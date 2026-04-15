package models

import "time"

type (
	SavingType string
	DailyLimit int
)

const (
	SavingTypeGroup    SavingType = "group"
	SavingTypePersonal SavingType = "personal"
	SavingTypeFlexible SavingType = "flexible"
	DailyLimit5        DailyLimit = 5
	DailyLimit10       DailyLimit = 10
)

type User struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Password   string     `json:"-"`
	SavingType SavingType `json:"saving_type"`
	DailyLimit DailyLimit `json:"daily_limit"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
