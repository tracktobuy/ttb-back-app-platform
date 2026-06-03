package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, item domain.User) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, item domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.CrudRepository[domain.User]
}

func NewUserService(repo repository.CrudRepository[domain.User]) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, item domain.User) (*domain.User, error) {
	return s.repo.Create(ctx, item)
}

func (s *userService) Get(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

func (s *userService) Update(ctx context.Context, id string, item domain.User) (*domain.User, error) {
	return nil, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return nil
}
