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
	GetStoreByUUID(ctx context.Context, storeUUID string) (response.Store, error)
	Delete(ctx context.Context, storeId primitive.ObjectID) error
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

	stores := s.formatStores(strs)
	return stores, nil
}

func (s *storeService) Delete(ctx context.Context, storeId primitive.ObjectID) error {
	return s.repo.Delete(ctx, storeId)
}

func (r *storeService) GetStoreByUUID(ctx context.Context, storeUUID string) (response.Store, error) {

	str, err := r.repo.GetStoreByUUID(ctx, storeUUID)
	if err != nil {
		return response.Store{}, err
	}

	store := r.formatStore(str)
	return store, nil
}

func (r *storeService) formatStores(input []domain.Store) []response.Store {
	var stores []response.Store

	for _, store := range input {
		tmp := r.formatStore(store)
		stores = append(stores, tmp)
	}

	return stores
}

func (r *storeService) formatStore(input domain.Store) response.Store {

	return response.Store{
		ID:           input.ID,
		UUID:         input.UUID,
		Price:        input.Price,
		ShippingCost: input.ShippingCost,
		Currency:     input.Currency,
		Domain:       input.Domain,
		Name:         input.Name,
		BestOption:   input.BestOption,
		URL:          input.URL,
		CreatedAt:    helper.DateTime(input.CreatedAt),
		UpdatedAt:    helper.DateTime(input.UpdatedAt),
	}
}
