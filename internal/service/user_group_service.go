package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type UserGroupService interface {
	Create(ctx context.Context, item domain.UserGroup) (*domain.UserGroup, error)
}

type userGroupService struct {
	repo repository.UserGroupRepository
}

func NewUserGroupService(repo repository.UserGroupRepository) UserGroupService {
	return &userGroupService{
		repo: repo,
	}
}

func (s *userGroupService) Create(ctx context.Context, item domain.UserGroup) (*domain.UserGroup, error) {
	return s.repo.Create(ctx, item)
}
