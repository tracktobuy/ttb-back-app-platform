package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type groupService struct {
	repo repository.CrudRepository[domain.Group]
}

func NewGroupService(repo repository.CrudRepository[domain.Group]) CrudService[domain.Group] {
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
