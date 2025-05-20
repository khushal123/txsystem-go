package repository

import (
	"context"
	"errors"
	"txsystem/internal/transaction/models"

	"gorm.io/gorm"
	// adjust import path accordingly
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByID(ctx context.Context, id uint) (*models.Transaction, error)
	Update(ctx context.Context, tx *models.Transaction) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]models.Transaction, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepo) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var tx models.Transaction
	result := r.db.WithContext(ctx).First(&tx, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tx, result.Error
}

func (r *transactionRepo) Update(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *transactionRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Transaction{}, id).Error
}

func (r *transactionRepo) List(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&transactions)
	return transactions, result.Error
}
