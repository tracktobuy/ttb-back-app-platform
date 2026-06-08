package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Create(ctx context.Context, item domain.User) (*domain.User, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, id string, item domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, item domain.User) (*domain.User, error) {
	return s.repo.Create(ctx, item)
}

func (s *userService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *userService) Update(ctx context.Context, id string, item domain.User) (*domain.User, error) {
	return nil, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return nil
}
