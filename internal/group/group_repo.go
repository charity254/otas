package group

import (
	"database/sql"
	"otas/models"
)

type GroupRepository struct {
	DB *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}
func (r *GroupRepository) GetGroupProgress(groupID int) (*models.GroupProgress, error) {
	query := `
		SELECT 
			id,
			name,
			target_amount - current_amount AS remaining_amount
		FROM groups
		WHERE id = $1
	`
	progress := &models.GroupProgress{}
	err := r.DB.QueryRow(query, groupID).Scan(
		&progress.GroupID,
		&progress.GroupName,
		&progress.RemainingAmount,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}
