package models

import (
	"time"
	"txsystem/pkg/common/types"
)

type Transaction struct {
	ID                 uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Amount             float64
	Currency           string
	Description        string
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	SourceAccount      string
	DestinationAccount string
	TransactionType    string
	Status             types.TransactionStatus
	TransactionID      string
}
