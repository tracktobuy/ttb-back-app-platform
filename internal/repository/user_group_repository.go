package repository

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGroupRepository interface {
	Create(ctx context.Context, item domain.UserGroup) (*domain.UserGroup, error)
}

type mongoUserGroupRepo struct {
	collection *mongo.Collection
}

func NewUserGroupRepo(db *mongo.Database) UserGroupRepository {
	return &mongoUserGroupRepo{
		collection: db.Collection(USER_GROUP_COLLECTION_NAME),
	}
}

func (r *mongoUserGroupRepo) Create(ctx context.Context, item domain.UserGroup) (*domain.UserGroup, error) {

	item.ID = primitive.NewObjectID()
	item.JoinedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
