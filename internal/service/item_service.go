package service

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemServiceInterface interface {
	CrudService[domain.Item]
	GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error)
}

type itemService struct {
	repo repository.ItemRepositoryInterface
}

func NewItemService(repo repository.ItemRepositoryInterface) ItemServiceInterface {
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

func (s *itemService) Get(ctx context.Context, id string) (*domain.Item, error) {
	return nil, nil
}

func (s *itemService) GetAll(ctx context.Context) ([]domain.Item, error) {
	return []domain.Item{}, nil
}

func (s *itemService) Update(ctx context.Context, id string, item domain.Item) (*domain.Item, error) {
	return nil, nil
}

func (s *itemService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *itemService) GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error) {

	items, err := s.repo.GetAllByGroupID(ctx, groupID)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}
