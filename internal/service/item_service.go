package service

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemService interface {
	Create(ctx context.Context, item domain.Item) (*response.Item, error)
	GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error)
}

type itemService struct {
	repo repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{
		repo: repo,
	}
}

func (s *itemService) Create(ctx context.Context, item domain.Item) (*response.Item, error) {

	item.ID = primitive.NewObjectID()
	item.UUID = helper.GenerateUUIDV7()
	item.Version = 1
	item.CreatedAt = time.Now().UTC()

	newItem, err := s.repo.Create(ctx, item)

	if err != nil {
		return nil, err
	}

	response := &response.Item{
		ID:          newItem.ID,
		UUID:        newItem.UUID,
		Title:       newItem.Title,
		Images:      newItem.Images,
		Labels:      newItem.Labels,
		CreatedAt:   helper.DateTime(newItem.CreatedAt),
		UpdatedAt:   helper.DateTime(newItem.UpdatedAt),
		CreatedByID: newItem.CreatedBy,
	}

	return response, nil
}

func (s *itemService) GetAllByGroupID(ctx context.Context, groupID primitive.ObjectID) ([]domain.Item, error) {

	items, err := s.repo.GetAllByGroupID(ctx, groupID)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}
