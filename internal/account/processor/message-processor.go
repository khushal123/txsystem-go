package processor

import (
	"txsystem/internal/account/service"
	"txsystem/pkg/common/types"

	"gorm.io/gorm"
)

type messageProcessor struct {
	acs *service.AccountService
}

func NewMessageProcessor(db *gorm.DB) types.MessageProcessor {
	return &messageProcessor{
		acs: service.NewAccountService(db),
	}
}

func (mp *messageProcessor) ProcessMessage(message string) error {
	return nil
}
