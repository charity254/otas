package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Deduction   float64   `json:"deduction"`
	AllocatedTo string    `json:"allocated_to"`
	CreatedAt   time.Time `json:"created_at"`
}
