package transaction

import (
	"database/sql"
	"errors"

	"otas/internal/account"
	"otas/internal/user"
	"otas/models"
)

type transactionService struct {
	repo        *transactionRepository
	accountRepo *account.AccountRepository
	userRepo    *user.UserRepository
}

func NewTransactionService(repo *transactionRepository, accountRepo *account.AccountRepository, userRepo *user.UserRepository) *transactionService {
	return &transactionService{repo: repo, accountRepo: accountRepo, userRepo: userRepo}
}

func (s *transactionService) ProcessTransaction(userID int, amount float64, user *models.User) (*models.Transaction, error) {
	// 1. Count today's transactions
	count, err := s.repo.CountTodayTransactions(userID)
	if err != nil {
		return nil, errors.New("failed to check daily limit")
	}

	// 2. Daily limit already hit — record transaction, no deduction
	if count >= int(user.DailyLimit) {
		return s.repo.CreateTransaction(&models.Transaction{
			UserID:      userID,
			Amount:      amount,
			Deduction:   0,
			Type:        models.TransactionTypeDeduction,
			AllocatedTo: "flexible",
		})
	}

	// 3. Calculate 10% deduction
	deduction := amount * 0.10

	// 4. Decide which account gets the deduction
	allocatedTo, err := s.getAllocatedAccount(user, count)
	if err != nil {
		return nil, err
	}

	// 5. Wrap balance update + record in DB transaction
	var savedTx *models.Transaction
	err = s.repo.WithTx(func(tx *sql.Tx) error {
		// update balance — transactional version
		if err := s.accountRepo.UpdateBalanceTx(tx, userID, allocatedTo, deduction); err != nil {
			return err
		}

		// record transaction — transactional version
		t, err := s.repo.InsertTransaction(tx, &models.Transaction{
			UserID:      userID,
			Amount:      amount,
			Deduction:   deduction,
			AllocatedTo: string(allocatedTo),
			Type:        models.TransactionTypeDeduction,
		})
		if err != nil {
			return err
		}

		savedTx = t
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to process transaction")
	}

	return savedTx, nil
}

func (s *transactionService) getAllocatedAccount(user *models.User, txCount int) (models.AccountType, error) {
	switch {
	// Group + locked:(10) first 5 → group, next 5 → locked
	case user.SavingType == models.SavingTypeGroup && user.DailyLimit == models.DailyLimit10:
		if txCount < 5 {
			return models.AccountTypeGroup, nil
		}
		return models.AccountTypeLocked, nil

	// Group only(5): all 5 → group
	case user.SavingType == models.SavingTypeGroup:
		return models.AccountTypeGroup, nil

	// Personal + locked(10); first 5 → locked, next 5 → flexible
	case user.SavingType == models.SavingTypePersonal && user.DailyLimit == models.DailyLimit10:
		if txCount < 5 {
			return models.AccountTypeLocked, nil
		}
		return models.AccountTypeFlexible, nil

	// Personal and locked(5) : all 5 → locked
	case user.SavingType == models.SavingTypePersonal:
		return models.AccountTypeLocked, nil

	// Personal flexible only(5  or 10 ): all → flexible
	case user.SavingType == models.SavingTypeFlexible:
		return models.AccountTypeFlexible, nil
	}

	return "", errors.New("invalid saving configuration")
}

func (s *transactionService) GetMemberContribution(userID int, accountType string) (float64, error) {
	total, err := s.repo.GetTotalContributed(userID, accountType)
	if err != nil {
		return 0, errors.New("failed to get member contribution")
	}
	return total, nil
}

func (s *transactionService) GetUser(userID int) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *transactionService) GetGroupTotalSavings(userID int) (float64, error) {
	balance, err := s.accountRepo.GetAccountBalance(userID, models.AccountTypeGroup)
	if err != nil {
		return 0, errors.New("failed to get group total savings")
	}
	return balance, nil
}

func (s *transactionService) LogWithdrawalPayout(userID int, amount float64, withdrawalRequestID int) (*models.Transaction, error) {
	var savedTx *models.Transaction

	err := s.repo.WithTx(func(tx *sql.Tx) error {
		if err := s.accountRepo.UpdateBalanceTx(tx, userID, models.AccountTypeGroup, -amount); err != nil {
			return err
		}
		if err := s.accountRepo.UpdateBalanceTx(tx, userID, models.AccountTypeFlexible, amount); err != nil {
			return err
		}

		t, err := s.repo.InsertTransaction(tx, &models.Transaction{
			UserID:      userID,
			Amount:      amount,
			Deduction:   0,
			AllocatedTo: "flexible",
			Type:        models.TransactionTypeWithdrawalPayout,
		})
		if err != nil {
			return err
		}

		savedTx = t
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to log withdrawal payout")
	}

	return savedTx, nil
}

func (s *transactionService) LogRepayment(userID int, amount float64) (*models.Transaction, error) {
	var savedTx *models.Transaction

	err := s.repo.WithTx(func(tx *sql.Tx) error {
		if err := s.accountRepo.UpdateBalanceTx(tx, userID, models.AccountTypeFlexible, -amount); err != nil {
			return err
		}
		if err := s.accountRepo.UpdateBalanceTx(tx, userID, models.AccountTypeGroup, amount); err != nil {
			return err
		}

		t, err := s.repo.InsertTransaction(tx, &models.Transaction{
			UserID:      userID,
			Amount:      amount,
			Deduction:   0,
			AllocatedTo: "group",
			Type:        models.TransactionTypeRepayment,
		})
		if err != nil {
			return err
		}

		savedTx = t
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to log repayment")
	}

	return savedTx, nil
}
