// internal/user/user_service.go
package user

import (
	"errors"

	"otas/models"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo *userRepository
}

func NewService(repo *userRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) (*models.User, error) {
	// 1. Check if email already exists
	existing, _ := s.repo.GetUserByEmail(user.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash the password BEFORE saving
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = string(hashed)

	// 3. Create user
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// 3. Save user, password is now a hash
	if err := s.createAccounts(createdUser); err != nil {
		return nil, err
	}

	return createdUser, nil
}
