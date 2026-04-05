package models

import "time"

type Group struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GroupProgress struct {
	GroupID          int     `json:"group_id"`
	GroupName        string  `json:"group_name"`
	TargetAmount     float64 `json:"target_amount"`
	CurrentAmount    float64 `json:"current_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	UserContribution float64 `json:"user_contribution"`
}
