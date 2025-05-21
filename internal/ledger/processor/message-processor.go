package processor

import (
	"context"
	"time"
	"txsystem/internal/common/types"
	"txsystem/internal/ledger/models"
	"txsystem/internal/ledger/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type messageProcessor struct {
	s *service.LedgerService
}

func NewMessageProcessor(db *gorm.DB) types.MessageProcessor {
	return &messageProcessor{}
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
