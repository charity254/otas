// internal/user/user_service.go
package user

import (
	"errors"

	"otas/internal/account"
	"otas/models"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo        *userRepository
	accountRepo *account.AccountRepository
}

func NewUserService(repo *userRepository, accountRepo *account.AccountRepository) *userService {
	return &userService{repo: repo, accountRepo: accountRepo}
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

	if err := s.createAccounts(createdUser); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) createAccounts(user *models.User) error {
	// always create main account
	if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeMain); err != nil {
		return errors.New("failed to create main account")
	}

	switch {
	case user.SavingType == models.SavingTypeGroup && user.DailyLimit == models.DailyLimit10:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeGroup); err != nil {
			return errors.New("failed to create group account")
		}
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}

	case user.SavingType == models.SavingTypeGroup && user.DailyLimit == models.DailyLimit5:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeGroup); err != nil {
			return errors.New("failed to create group account")
		}

	case user.SavingType == models.SavingTypePersonal && user.DailyLimit == models.DailyLimit10:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeFlexible); err != nil {
			return errors.New("failed to create flexible account")
		}

	case user.SavingType == models.SavingTypePersonal && user.DailyLimit == models.DailyLimit5:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}
	}

	return nil
}

func (s *userService) Login(email, password string) (*models.User, error) {
	//verify user email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}
