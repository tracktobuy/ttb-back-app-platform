package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemRepository interface {
	Create(ctx context.Context, item domain.Item) (*domain.Item, error)
	Get(ctx context.Context, id string) (*domain.Item, error)
	Delete(ctx context.Context, id string) error
	GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error)
}

type itemRepo struct {
	collection *mongo.Collection
}

func NewItemRepository(db *mongo.Database) ItemRepository {
	return &itemRepo{collection: db.Collection("items")}
}

func (r *itemRepo) Create(ctx context.Context, item domain.Item) (*domain.Item, error) {
	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *itemRepo) Get(ctx context.Context, id string) (*domain.Item, error) {
	var item domain.Item
	err := r.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) Delete(ctx context.Context, id string) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"uuid": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *itemRepo) GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error) {
	var items []domain.Item

	filterList := []primitive.ObjectID{groupID}

	filter := bson.M{
		"$in": filterList,
	}

	cursor, err := r.collection.Find(ctx, bson.M{"groups": filter})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}
