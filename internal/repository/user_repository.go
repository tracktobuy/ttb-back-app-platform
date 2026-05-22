package repository

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) CrudRepository[domain.User] {
	return &mongoUserRepo{
		collection: db.Collection("users"),
	}
}

func (u *mongoUserRepo) Create(ctx context.Context, item domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (u *mongoUserRepo) Get(ctx context.Context, id string) (domain.User, error) {
	return domain.User{}, nil
}

func (u *mongoUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	return []domain.User{}, nil
}

func (u *mongoUserRepo) Update(ctx context.Context, item domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (u *mongoUserRepo) Delete(ctx context.Context, id string) error {
	return nil
}
