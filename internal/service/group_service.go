package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type GroupService struct {
	repo repository.CrudRepository[domain.Group]
}

func NewGroupService(repo repository.CrudRepository[domain.Group]) CrudService[domain.Group] {
	return &GroupService{
		repo: repo,
	}
}

func NewGroupServiceImplementation(repo repository.CrudRepository[domain.Group]) GroupService {
	return GroupService{
		repo: repo,
	}
}

func (s *GroupService) Create(ctx context.Context, item domain.Group) (*domain.Group, error) {
	return s.repo.Create(ctx, item)
}

func (s *GroupService) Get(ctx context.Context, id string) (*domain.Group, error) {
	return s.repo.Get(ctx, id)
}

func (s *GroupService) GetAll(ctx context.Context) ([]domain.Group, error) {
	return s.repo.GetAll(ctx)
}

func (s *GroupService) Update(ctx context.Context, id string, item domain.Group) (*domain.Group, error) {
	return nil, nil
}

func (s *GroupService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *GroupService) CreateDefaultGroup(ctx context.Context, user domain.User) (*domain.Group, error) {

	value, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	defaultGroup := domain.Group{
		UUID:           value.String(),
		Name:           "Wishlist",
		Budget:         0.0,
		BudgetCurrency: "BRL",
		CreatedBy:      user.ID,
	}

	return s.repo.Create(ctx, defaultGroup)
}
