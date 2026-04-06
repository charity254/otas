package group

import (
	"errors"
	"math"

	"otas/models"
)

type groupService struct {
	repo *GroupRepository
}

func NewGroupService(repo *GroupRepository) *groupService {
	return &groupService{repo: repo}
}
func (s *groupService) CreateGroup(name string, targetAmount float64, createdBy int) (*models.Group, error) {
	// 1. Validate target amount
	if targetAmount <= 0 {
		return nil, errors.New("target amount must be greater than zero")
	}

	// 2. Create group in DB
	group, err := s.repo.CreateGroup(name, targetAmount, createdBy)
	if err != nil {
		return nil, errors.New("failed to create group")
	}

	// 3. Assign creator as first member
	if err := s.repo.AddGroupMember(group.ID, createdBy); err != nil {
		return nil, errors.New("failed to assign creator as group member")
	}

	// 4. creator is first member — contribution = full target
	if err := s.repo.UpdateAllMemberContributions(group.ID, targetAmount); err != nil {
		return nil, errors.New("failed to set initial contribution")
	}

	return group, nil
}

func (s *groupService) GetGroupProgress(groupID, userID int) (*models.GroupProgress, error) {
	// 1. Get group progress — target, remaining
	progress, err := s.repo.GetGroupProgress(groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	// 2. Get this user's individual contribution
	contribution, err := s.repo.GetMemberContribution(groupID, userID)
	if err != nil {
		return nil, errors.New("failed to get member contribution")
	}

	// 3. Attach contribution to progress
	progress.UserContribution = contribution

	return progress, nil
}

func (s *groupService) GetMemberContribution(groupID, userID int) (float64, error) {
	contribution, err := s.repo.GetMemberContribution(groupID, userID)
	if err != nil {
		return 0, errors.New("failed to get member contribution")
	}
	return contribution, nil
}

func (s *groupService) JoinGroup(groupID, userID int) error {
	// 1. Check if already a member
	isMember, err := s.repo.IsGroupMember(groupID, userID)
	if err != nil {
		return errors.New("failed to check membership")
	}
	if isMember {
		return errors.New("already a member of this group")
	}

	// 2. Add user to group
	if err := s.repo.AddGroupMember(groupID, userID); err != nil {
		return errors.New("failed to join group")
	}

	// 3. Recalculate contribution per member now that count changed
	contribution, err := s.CalculateMemberContribution(groupID)
	if err != nil {
		return errors.New("failed to recalculate contributions")
	}

	// 4. Update all members' contribution amount
	if err := s.repo.UpdateAllMemberContributions(groupID, contribution); err != nil {
		return errors.New("failed to update member contributions")
	}

	return nil
}

func (s *groupService) CalculateMemberContribution(groupID int) (float64, error) {
	// 1. Get group target
	target, err := s.repo.GetGroupTarget(groupID)
	if err != nil {
		return 0, errors.New("failed to get group target")
	}

	// 2. Get member count
	memberCount, err := s.repo.GetGroupMemberCount(groupID)
	if err != nil {
		return 0, errors.New("failed to get member count")
	}

	if target == 0 {
		return 0, errors.New("group has not set a target")
	}

	if memberCount == 0 {
		return 0, errors.New("group has no members")
	}

	// 3. Calculate contribution per member
	// math.Round to 2 decimal places to handle rounding
	contribution := math.Round((target/float64(memberCount))*100) / 100

	return contribution, nil
}
