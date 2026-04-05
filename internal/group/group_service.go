package group

import (
	"errors"

	"otas/models"
)

type groupService struct {
	repo *GroupRepository
}

func NewGroupService(repo *GroupRepository) *groupService {
	return &groupService{repo: repo}
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
