package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
)

type GroupService interface {
	Create(ctx context.Context, item request.Group, user domain.User) (*response.Group, error)
	Get(ctx context.Context, id string) (*response.Group, error)
	Update(ctx context.Context, id string, item request.Group) (*response.Group, error)
	CreateDefaultGroup(ctx context.Context, user domain.User) (*response.Group, error)
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{
		repo: repo,
	}
}

func (s *groupService) Create(ctx context.Context, reqGroup request.Group, user domain.User) (*response.Group, error) {

	item := domain.Group{
		Name:           reqGroup.Name,
		Budget:         reqGroup.Budget,
		BudgetCurrency: reqGroup.BudgetCurrency,
		CreatedBy:      user.ID,
	}

	newGroup, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return s.formatGroupResponse(newGroup), nil
}

func (s *groupService) Get(ctx context.Context, id string) (*response.Group, error) {
	group, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.formatGroupResponse(group), nil
}

func (s *groupService) Update(ctx context.Context, id string, reqGroup request.Group) (*response.Group, error) {

	grp, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	grp.Name = reqGroup.Name
	grp.Budget = reqGroup.Budget
	grp.BudgetCurrency = reqGroup.BudgetCurrency

	group, err := s.repo.Update(ctx, grp)
	if err != nil {
		return nil, err
	}

	return s.formatGroupResponse(group), err
}

func (s *groupService) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *groupService) CreateDefaultGroup(ctx context.Context, user domain.User) (*response.Group, error) {

	defaultGroup := domain.Group{
		Name:           user.Name + "'s Wishlist",
		Budget:         0.0,
		BudgetCurrency: "BRL",
		CreatedBy:      user.ID,
	}

	newGroup, err := s.repo.Create(ctx, defaultGroup)
	if err != nil {
		return nil, err
	}

	return s.formatGroupResponse(newGroup), nil
}

func (s *groupService) formatGroupResponse(group *domain.Group) *response.Group {
	return &response.Group{
		ID:             group.ID,
		UUID:           group.UUID,
		Name:           group.Name,
		Budget:         group.Budget,
		BudgetCurrency: group.BudgetCurrency,
		CreatedAt:      helper.DateTime(group.CreatedAt),
		UpdatedAt:      helper.DateTime(group.UpdatedAt),
	}
}
