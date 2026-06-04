package service

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemService interface {
	Create(ctx context.Context, item domain.Item) (*domain.Item, error)
	GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error)
}

type itemService struct {
	repo repository.ItemRepositoryInterface
}

func NewItemService(repo repository.ItemRepositoryInterface) ItemService {
	return &itemService{
		repo: repo,
	}
}

func (s *itemService) Create(ctx context.Context, item domain.Item) (*domain.Item, error) {

	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now().UTC()

	newItem, err := s.repo.Create(ctx, item)

	if err != nil {
		return nil, err
	}

	return newItem, nil
}

func (s *itemService) GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error) {

	items, err := s.repo.GetAllByGroupID(ctx, groupID)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}
