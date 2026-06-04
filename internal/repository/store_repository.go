package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreRepository interface {
	Create(ctx context.Context, item domain.Store) (*domain.Store, error)
}

type storeRepo struct {
	collection *mongo.Collection
}

func NewStoreRepository(db *mongo.Database) StoreRepository {
	return &storeRepo{collection: db.Collection("stores")}
}

func (r *storeRepo) Create(ctx context.Context, item domain.Store) (*domain.Store, error) {

	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
