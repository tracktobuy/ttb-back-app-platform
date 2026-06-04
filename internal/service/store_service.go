package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type StoreService interface {
	Create(ctx context.Context, item domain.Store) (*domain.Store, error)
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{
		repo: repo,
	}
}

func (s *storeService) Create(ctx context.Context, item domain.Store) (*domain.Store, error) {
	item.UUID = helper.GenerateUUIDV7()
	return s.repo.Create(ctx, item)
}
