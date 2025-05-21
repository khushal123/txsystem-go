package service

import (
	"context"
	"time"
	"txsystem/internal/ledger/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LedgerService struct {
	collection *mongo.Collection
}

func NewLedgerService(db *mongo.Database) *LedgerService {
	return &LedgerService{
		collection: db.Collection("ledger"),
	}
}

func (s *LedgerService) CreateLedger(ctx context.Context, ledger *models.Ledger) error {
	ledger.CreatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, ledger)
	if err != nil {
		return err
	}
	return nil
}

func (s *LedgerService) ListLedgers(ctx context.Context) ([]*models.Ledger, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ledgers []*models.Ledger
	for cursor.Next(ctx) {
		var l models.Ledger
		if err := cursor.Decode(&l); err != nil {
			return nil, err
		}
		ledgers = append(ledgers, &l)
	}
	return ledgers, cursor.Err()
}
