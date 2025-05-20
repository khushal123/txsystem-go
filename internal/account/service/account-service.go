package service

import (
	"context"
	"fmt"
	"txsystem/internal/account/models"

	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (as *AccountService) CreateAccount(ctx context.Context, owner, currency string, initialBalance float64) (*models.Account, error) {
	account := &models.Account{
		Owner:    owner,
		Currency: currency,
		Balance:  initialBalance,
	}
	if err := as.db.WithContext(ctx).Create(account).Error; err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return account, nil
}

func (as *AccountService) GetAccount(ctx context.Context, id uint) (*models.Account, error) {
	var account models.Account
	if err := as.db.WithContext(ctx).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (as *AccountService) TransferBalance(ctx context.Context, fromID, toID uint, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	if fromID == toID {
		return fmt.Errorf("source and destination accounts cannot be the same")
	}

	tx := as.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var fromAccount models.Account
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&fromAccount, fromID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get source account: %w", err)
	}

	// Check for sufficient balance
	if fromAccount.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient balance in source account")
	}

	var toAccount models.Account
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&toAccount, toID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get destination account: %w", err)
	}

	fromAccount.Balance -= amount
	if err := tx.Save(&fromAccount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update source account: %w", err)
	}

	toAccount.Balance += amount
	if err := tx.Save(&toAccount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
