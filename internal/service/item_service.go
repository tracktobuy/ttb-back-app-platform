package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type ItemService interface {
	Create(ctx context.Context, user *domain.User, item domain.Item) (*response.Item, error)
	GetByUUID(ctx context.Context, itemUUID string) (*response.Item, error)
	GetAllByUserId(ctx context.Context, user *domain.User) ([]response.Item, error)
}

type itemService struct {
	log  logger.Logger
	repo repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) ItemService {

	log := logger.NewLogger()
	log.SetServiceName("ItemService")

	return &itemService{
		repo: repo,
		log:  log,
	}
}

func (s *itemService) Create(ctx context.Context, user *domain.User, item domain.Item) (*response.Item, error) {

	newItem, err := s.repo.Create(ctx, item)

	if err != nil {
		return nil, err
	}

	response := s.formatItemResponse(user, *newItem)

	return &response, nil
}

func (s *itemService) GetByUUID(ctx context.Context, itemUUID string) (*response.Item, error) {

	it, err := s.repo.Get(ctx, itemUUID)
	if err != nil {
		return nil, err
	}

	item := s.formatItemResponse(nil, *it)
	return &item, nil
}

func (s *itemService) GetAllByUserId(ctx context.Context, user *domain.User) ([]response.Item, error) {

	items, err := s.repo.GetAllByUserId(ctx, user.ID)
	if err != nil {
		return []response.Item{}, err
	}
	resp := s.formatItemsResponse(user, items)
	return resp, nil
}

func (s *itemService) formatItemsResponse(user *domain.User, items []domain.Item) []response.Item {
	var resp []response.Item
	for _, item := range items {
		tmp := s.formatItemResponse(user, item)
		resp = append(resp, tmp)
	}
	return resp
}

func (s *itemService) formatItemResponse(user *domain.User, item domain.Item) response.Item {

	i := response.Item{
		ID:        item.ID,
		UUID:      item.UUID,
		Title:     item.Title,
		Images:    item.Images,
		CreatedAt: helper.DateTime(item.CreatedAt),
		UpdatedAt: helper.DateTime(item.UpdatedAt),
		Labels:    item.Labels,
		Stores:    []string{},
		Groups:    []string{},
	}

	if user != nil {
		i.CreatedBy = response.User{
			ID:       user.ID,
			UUID:     user.UUID,
			Name:     user.Name,
			Username: user.Username,
		}
	}

	return i
}
