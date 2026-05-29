package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type GroupServiceInterface interface {
	CrudService[domain.Group]
	CreateDefaultGroup(ctx context.Context, user domain.User) (*domain.Group, error)
}

type groupService struct {
	repo repository.CrudRepository[domain.Group]
}

func NewGroupService(repo repository.CrudRepository[domain.Group]) GroupServiceInterface {
	return &groupService{
		repo: repo,
	}
}

func NewGroupServiceImplementation(repo repository.CrudRepository[domain.Group]) GroupServiceInterface {
	return &groupService{
		repo: repo,
	}
}

func (s *groupService) Create(ctx context.Context, item domain.Group) (*domain.Group, error) {
	return s.repo.Create(ctx, item)
}

func (s *groupService) Get(ctx context.Context, id string) (*domain.Group, error) {
	return s.repo.Get(ctx, id)
}

func (s *groupService) GetAll(ctx context.Context) ([]domain.Group, error) {
	return s.repo.GetAll(ctx)
}

func (s *groupService) Update(ctx context.Context, id string, item domain.Group) (*domain.Group, error) {
	return nil, nil
}

func (s *groupService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *groupService) CreateDefaultGroup(ctx context.Context, user domain.User) (*domain.Group, error) {

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
