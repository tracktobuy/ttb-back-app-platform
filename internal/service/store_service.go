package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreService interface {
	Create(ctx context.Context, item domain.Store) (*domain.Store, error)
	GetStoresByItemId(ctx context.Context, itemId primitive.ObjectID) ([]response.Store, error)
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

func (s *storeService) GetStoresByItemId(ctx context.Context, itemId primitive.ObjectID) ([]response.Store, error) {

	strs, err := s.repo.GetStoresByItemId(ctx, itemId)

	if err != nil {
		return []response.Store{}, err
	}

	var stores []response.Store

	for _, str := range strs {
		tmp := response.Store{
			ID:           str.ID,
			UUID:         str.UUID,
			Price:        str.Price,
			ShippingCost: str.ShippingCost,
			Currency:     str.Currency,
			Domain:       str.Domain,
			Name:         str.Name,
			BestOption:   str.BestOption,
			URL:          str.URL,
			CreatedAt:    helper.DateTime(str.CreatedAt),
			UpdatedAt:    helper.DateTime(str.UpdatedAt),
		}

		stores = append(stores, tmp)
	}

	return stores, nil
}
