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
		SELECT id,name, current_amount, target_amount - current_amount AS remaining_amount
		FROM groups
		WHERE id = $1
	`
	progress := &models.GroupProgress{}
	err := r.DB.QueryRow(query, groupID).Scan(
		&progress.GroupID,
		&progress.GroupName,
		&progress.CurrentAmount,
		&progress.RemainingAmount,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *GroupRepository) GetMemberContribution(groupID, userID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(t.deduction), 0)
		FROM transactions t
		JOIN group_members gm ON gm.user_id = t.user_id
		WHERE t.user_id = $1
		AND t.allocated_to = 'group'
		AND gm.group_id = $2
	`
	var total float64
	err := r.DB.QueryRow(query, userID, groupID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *GroupRepository) IsGroupMember(groupID, userID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM group_members
			WHERE group_id = $1
			AND user_id = $2
		)
	`
	var exists bool
	err := r.DB.QueryRow(query, groupID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *GroupRepository) GetGroupMemberCount(groupID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM group_members
		WHERE group_id = $1
	`
	var count int
	err := r.DB.QueryRow(query, groupID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupRepository) GetGroupTarget(groupID int) (float64, error) {
	query := `
		SELECT target_amount FROM groups
		WHERE id = $1
	`
	var target float64
	err := r.DB.QueryRow(query, groupID).Scan(&target)
	if err != nil {
		return 0, err
	}
	return target, nil
}

func (r *GroupRepository) UpdateAllMemberContributions(groupID int, contribution float64) error {
	query := `
		UPDATE group_members
		SET contribution_amount = $1
		WHERE group_id = $2
	`
	_, err := r.DB.Exec(query, contribution, groupID)
	return err
}
