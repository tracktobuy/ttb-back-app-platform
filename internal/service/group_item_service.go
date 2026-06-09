package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type GroupItemService interface {
	Create(ctx context.Context, item domain.GroupItem) (*domain.GroupItem, error)
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
