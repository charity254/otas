package user

import (
	"errors"

	"otas/internal/account"
	"otas/models"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo        *UserRepository
	accountRepo *account.AccountRepository
}

func NewUserService(repo *UserRepository, accountRepo *account.AccountRepository) *userService {
	return &userService{repo: repo, accountRepo: accountRepo}
}

func (s *userService) Register(user *models.User) (*models.User, error) {
	// 1. Check if email already exists
	existing, _ := s.repo.GetUserByEmail(user.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash password
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

	// 4. Create correct accounts based on saving config
	if err := s.createAccounts(createdUser); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) Login(email, password string) (*models.User, error) {
	// 1. Find user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *userService) createAccounts(user *models.User) error {
	// always create main account for every user
	if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeMain); err != nil {
		return errors.New("failed to create main account")
	}

	switch {
	// Group + 10: group + locked
	case user.SavingType == models.SavingTypeGroup && user.DailyLimit == models.DailyLimit10:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeGroup); err != nil {
			return errors.New("failed to create group account")
		}
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}

	// Group + 5: group only
	case user.SavingType == models.SavingTypeGroup && user.DailyLimit == models.DailyLimit5:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeGroup); err != nil {
			return errors.New("failed to create group account")
		}

	// Personal + locked + 10: locked + flexible
	case user.SavingType == models.SavingTypePersonal && user.DailyLimit == models.DailyLimit10:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeFlexible); err != nil {
			return errors.New("failed to create flexible account")
		}

	// Personal + locked + 5: locked only
	case user.SavingType == models.SavingTypePersonal && user.DailyLimit == models.DailyLimit5:
		if _, err := s.accountRepo.CreateAccount(user.ID, models.AccountTypeLocked); err != nil {
			return errors.New("failed to create locked account")
		}

	// Personal flexible only: main account serves as flexible, nothing extra needed
	case user.SavingType == models.SavingTypeFlexible:
		// main account already created above
	}

	return nil
}
