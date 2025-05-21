package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ledger struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount    float64            `bson:"amount" json:"amount"`
	AccountID string             `bson:"account_id" json:"account_id"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
