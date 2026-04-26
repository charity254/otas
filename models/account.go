package models

import "time"

type AccountType string

const (
	AccountTypeFlexible AccountType = "flexible"
	AccountTypeLocked   AccountType = "locked"
	AccountTypeGroup    AccountType = "group"
)

type Account struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Type      AccountType `json:"type"`
	Balance   float64     `json:"balance"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
