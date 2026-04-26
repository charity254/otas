package models

import "time"

type TransactionType string

const (
	TransactionTypeDeduction        TransactionType = "deduction"
	TransactionTypeWithdrawalPayout TransactionType = "withdrawal_payout"
	TransactionTypeRepayment        TransactionType = "repayment"
)

type Transaction struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	Amount      float64         `json:"amount"`
	Deduction   float64         `json:"deduction"`
	AllocatedTo string          `json:"allocated_to"`
	Type        TransactionType `json:"type"`
	CreatedAt   time.Time       `json:"created_at"`
}
