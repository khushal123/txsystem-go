package service

import (
	"context"
	"encoding/json"
	"fmt"
	"txsystem/internal/common/types"
	"txsystem/internal/transaction/models"
	"txsystem/internal/transaction/repository"
)

type TransactionService struct {
	kc   types.ProducerConnection
	repo repository.TransactionRepository
}

func NewTransactionService(kc types.ProducerConnection, repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		kc:   kc,
		repo: repo,
	}
}

// toTransactionModel maps a request DTO to the persistence model.
func toTransactionModel(req *types.TransactionRequest) *models.Transaction {
	return &models.Transaction{
		Amount:             req.Amount,
		Description:        req.Description,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		TransactionType:    req.TransactionType,
		Status:             types.StatusPending,
		Currency:           "USD",
	}
}

// toTransactionResponse maps a persistence model to the response DTO.
func toTransactionResponse(m *models.Transaction) *types.TransactionResponse {
	return &types.TransactionResponse{
		ID:                 uint64(m.ID),
		Amount:             m.Amount,
		Description:        m.Description,
		SourceAccount:      m.SourceAccount,
		DestinationAccount: m.DestinationAccount,
		TransactionType:    m.TransactionType,
		Status:             string(m.Status),
		CreatedAt:          m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateTransaction persists a new transaction and sends an event.
func (ts *TransactionService) CreateTransaction(
	ctx context.Context,
	req *types.TransactionRequest,
) error {
	model := toTransactionModel(req)
	if err := ts.repo.Create(ctx, model); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	resp := toTransactionResponse(model)
	payload, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction response: %w", err)
	}

	if err := ts.kc.Produce(string(payload)); err != nil {
		return fmt.Errorf("failed to produce kafka event: %w", err)
	}

	return nil
}

// GetTransactions retrieves paginated transactions and maps them to response DTOs.
func (ts *TransactionService) GetTransactions(
	ctx context.Context,
) ([]*types.TransactionResponse, error) {
	modelsList, err := ts.repo.List(ctx, 100, 0)
	if err != nil {
		return nil, err
	}

	var respList []*types.TransactionResponse
	for _, m := range modelsList {
		respList = append(respList, toTransactionResponse(&m))
	}
	return respList, nil
}

// GetTransaction retrieves a single transaction by ID and maps it to a response DTO.
func (ts *TransactionService) GetTransaction(
	ctx context.Context,
	id uint,
) (*types.TransactionResponse, error) {
	m, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return toTransactionResponse(m), nil
}
