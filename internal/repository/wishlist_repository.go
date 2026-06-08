package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const WISHLIST_COLLECTION_NAME = "wishlists"

type WishlistRepository interface {
	Create(ctx context.Context, item domain.Wishlist) (*domain.Wishlist, error)
}

type mongoWishlistRepo struct {
	collection *mongo.Collection
}

func NewWishlistRepo(db *mongo.Database) WishlistRepository {
	return &mongoWishlistRepo{
		collection: db.Collection(WISHLIST_COLLECTION_NAME),
	}
}

func (r *mongoWishlistRepo) Create(ctx context.Context, item domain.Wishlist) (*domain.Wishlist, error) {

	item.ID = primitive.NewObjectID()
	item.JoinedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
