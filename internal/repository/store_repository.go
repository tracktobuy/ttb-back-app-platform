package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/mongo"
)

type storeRepo struct {
	collection *mongo.Collection
}

func NewStoreRepository(db *mongo.Database) CrudRepository[domain.Store] {
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

func (r *storeRepo) Get(ctx context.Context, id string) (*domain.Store, error) {
	return nil, nil
}

func (r *storeRepo) GetAll(ctx context.Context) ([]domain.Store, error) {
	return nil, nil
}

func (r *storeRepo) Update(ctx context.Context, item domain.Store) (*domain.Store, error) {
	return nil, nil
}

func (r *storeRepo) Delete(ctx context.Context, id string) error {
	return nil
}
