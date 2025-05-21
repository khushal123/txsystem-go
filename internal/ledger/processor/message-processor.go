package processor

import (
	"context"
	"time"
	"txsystem/internal/ledger/models"
	"txsystem/internal/ledger/service"
	"txsystem/pkg/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageProcessor struct {
	s *service.LedgerService
}

func NewMessageProcessor(db *mongo.Database) types.MessageProcessor {
	return &messageProcessor{
		s: service.NewLedgerService(db),
	}
}

func (mp *messageProcessor) ProcessMessage(message string) error {
	ctx := context.Background()
	ledger := &models.Ledger{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
	}
	err := mp.s.CreateLedger(ctx, ledger)
	if err != nil {
		return err
	}
	return nil
}
