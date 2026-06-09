package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const WISHLIST_COLLECTION_NAME = "wishlists"

type GroupItemRepository interface {
	Create(ctx context.Context, item domain.GroupItem) (*domain.GroupItem, error)
}

type mongoGroupItemRepo struct {
	collection *mongo.Collection
}

func NewGroupItemRepo(db *mongo.Database) GroupItemRepository {
	return &mongoGroupItemRepo{
		collection: db.Collection(WISHLIST_COLLECTION_NAME),
	}
}

func (r *mongoGroupItemRepo) Create(ctx context.Context, item domain.GroupItem) (*domain.GroupItem, error) {

	item.ID = primitive.NewObjectID()
	item.JoinedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
