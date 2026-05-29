package service

import (
	"context"
	"log/slog"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
)

type accountService struct {
	userService  CrudService[domain.User]
	groupService GroupServiceInterface
	ctx          context.Context
}

func NewAccountService(ctx context.Context, userService CrudService[domain.User], groupService GroupServiceInterface) *accountService {
	return &accountService{
		userService:  userService,
		groupService: groupService,
		ctx:          ctx,
	}
}

func (s *accountService) CreateAccount(user domain.User) (*domain.User, *domain.Group, error) {

	newUser, err := s.userService.Create(s.ctx, user)
	if err != nil {
		slog.Error("Error when creating new user account", "error", err.Error())
		return nil, nil, err
	}

	newGroup, err := s.groupService.CreateDefaultGroup(s.ctx, *newUser)
	if err != nil {
		slog.Error("Error when creating default group to user", "error", err.Error())
		return nil, nil, err
	}

	return newUser, newGroup, nil
}
