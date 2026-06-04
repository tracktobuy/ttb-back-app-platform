package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type StoreService interface {
	Create(ctx context.Context, item domain.Store) (*domain.Store, error)
	Get(ctx context.Context, id string) (*domain.Store, error)
	GetAll(ctx context.Context) ([]domain.Store, error)
	Update(ctx context.Context, id string, item domain.Store) (*domain.Store, error)
	Delete(ctx context.Context, id string) error
}

type storeService struct {
	repo repository.CrudRepository[domain.Store]
}

func NewStoreService(repo repository.CrudRepository[domain.Store]) StoreService {
	return &storeService{
		repo: repo,
	}
}

func (s *storeService) Create(ctx context.Context, item domain.Store) (*domain.Store, error) {
	item.UUID = helper.GenerateUUIDV7()
	return s.repo.Create(ctx, item)
}

func (s *storeService) Get(ctx context.Context, id string) (*domain.Store, error) {
	return s.repo.Get(ctx, id)
}

func (s *storeService) GetAll(ctx context.Context) ([]domain.Store, error) {
	return s.repo.GetAll(ctx)
}

func (s *storeService) Update(ctx context.Context, id string, item domain.Store) (*domain.Store, error) {
	return s.repo.Update(ctx, item)
}

func (s *storeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
