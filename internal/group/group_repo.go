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
func (r *GroupRepository) CreateGroup(name string, targetAmount float64, createdBy int) (*models.Group, error) {
	query := `
		INSERT INTO groups (name, target_amount, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, name, target_amount, current_amount, created_by, created_at, updated_at
	`
	group := &models.Group{}
	err := r.DB.QueryRow(query, name, targetAmount, createdBy).Scan(
		&group.ID,
		&group.Name,
		&group.TargetAmount,
		&group.CurrentAmount,
		&group.CreatedBy,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) AddGroupMember(groupID, userID int) error {
	query := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
	`
	_, err := r.DB.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) GetGroupProgress(groupID int) (*models.GroupProgress, error) {
	query := `
		SELECT id,name,target_amount - current_amount AS remaining_amount
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

// group_repo.go
func (r *GroupRepository) GetMemberContribution(groupID, userID int) (float64, error) {
	query := `
		SELECT contribution_amount FROM group_members
		WHERE group_id = $1
		AND user_id = $2
	`
	var contribution float64
	err := r.DB.QueryRow(query, groupID, userID).Scan(&contribution)
	if err != nil {
		return 0, err
	}
	return contribution, nil
}
