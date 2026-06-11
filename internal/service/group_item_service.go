package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupItemService interface {
	Create(ctx context.Context, item domain.GroupItem) (*domain.GroupItem, error)
	GetAllByGroupId(ctx context.Context, groupId primitive.ObjectID) ([]domain.GroupItem, error)
}

type groupItemService struct {
	repo repository.GroupItemRepository
}

func NewGroupItemService(repo repository.GroupItemRepository) GroupItemService {
	return &groupItemService{
		repo: repo,
	}
}

func (s *groupItemService) Create(ctx context.Context, item domain.GroupItem) (*domain.GroupItem, error) {
	return s.repo.Create(ctx, item)
}

func (s *groupItemService) GetAllByGroupId(ctx context.Context, groupId primitive.ObjectID) ([]domain.GroupItem, error) {
	return []domain.GroupItem{}, nil
}
