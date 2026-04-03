package models

import "time"

type WithdrawalStatus string

const (
	WithdrawalStatusPending   WithdrawalStatus = "pending"
	WithdrawalStatusApproved  WithdrawalStatus = "approved"
	WithdrawalStatusRejected  WithdrawalStatus = "rejected"
	WithdrawalStatusCancelled WithdrawalStatus = "cancelled"
)

type WithdrawalRequest struct {
	ID        int              `json:"id"`
	GroupID   int              `json:"group_id"`
	UserID    int              `json:"user_id"` // requester
	Amount    float64          `json:"amount"`
	Status    WithdrawalStatus `json:"status"`
	ExpiresAt time.Time        `json:"expires_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type WithdrawalApproval struct {
	ID                  int       `json:"id"`
	WithdrawalRequestID int       `json:"withdrawal_request_id"`
	UserID              int       `json:"user_id"`
	ContactInput        string    `json:"contact_input"`
	Status              string    `json:"status"` // 'approved', 'rejected'
	CreatedAt           time.Time `json:"created_at"`
}
