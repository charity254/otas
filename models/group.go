package models

import "time"

type Group struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	TargetAmount float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GroupMember struct {
	ID                 int       `json:"id"`
	GroupID            int       `json:"group_id"`
	UserID             int       `json:"user_id"`
	ContributionAmount float64   `json:"contribution_amount"`
	DebtAmount         float64   `json:"debt_amount"`
	JoinedAt           time.Time `json:"joined_at"`
}

type GroupProgress struct {
	GroupID          int     `json:"group_id"`
	GroupName        string  `json:"group_name"`
	TargetAmount     float64 `json:"target_amount"`
	CurrentAmount    float64 `json:"current_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	UserContribution float64 `json:"user_contribution"`
}
