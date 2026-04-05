package models

type GroupProgress struct {
	GroupID          int     `json:"group_id"`
	GroupName        string  `json:"group_name"`
	TargetAmount     float64 `json:"target_amount"`
	CurrentAmount    float64 `json:"current_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	UserContribution float64 `json:"user_contribution"`
}
