package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreRepository interface {
	Create(ctx context.Context, item domain.Store) (*domain.Store, error)
	GetStoresByItemId(ctx context.Context, itemId primitive.ObjectID) ([]domain.Store, error)
	GetStoreByUUID(ctx context.Context, storeUUID string) (domain.Store, error)
	Update(ctx context.Context, store domain.Store) (*domain.Store, error)
	Delete(ctx context.Context, storeId primitive.ObjectID) error
}

type storeRepo struct {
	log        logger.Logger
	collection *mongo.Collection
}

func NewStoreRepo(db *mongo.Database) StoreRepository {
	log := logger.NewLogger()
	log.SetRepositoryName("StoreRepository")

	return &storeRepo{
		collection: db.Collection(STORES_COLLECTION_NAME),
		log:        log,
	}
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

func (r *storeRepo) GetStoresByItemId(ctx context.Context, itemId primitive.ObjectID) ([]domain.Store, error) {

	r.log.Info("get stores by item id", "itemID", itemId)

	var stores []domain.Store

	cursor, err := r.collection.Find(ctx, bson.M{"itemId": itemId})

	if err != nil {
		return []domain.Store{}, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &stores); err != nil {
		return []domain.Store{}, err
	}

	return stores, nil
}

func (r *storeRepo) Update(ctx context.Context, store domain.Store) (*domain.Store, error) {

	version := store.Version
	store.UpdatedAt = time.Now().UTC()
	store.Version += 1

	updateData := bson.M{
		"$set": bson.M{
			"shippingCost": store.ShippingCost,
			"currency":     store.Currency,
			"version":      store.Version,
			"updatedAt":    store.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": store.ID, "version": version}, updateData)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *storeRepo) Delete(ctx context.Context, storeId primitive.ObjectID) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": storeId})

	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepo) GetStoreByUUID(ctx context.Context, storeUUID string) (domain.Store, error) {

	var str domain.Store
	if err := r.collection.FindOne(ctx, bson.M{"uuid": storeUUID}).Decode(&str); err != nil {
		return domain.Store{}, err
	}

	return str, nil
}
