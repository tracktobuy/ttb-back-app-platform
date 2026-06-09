package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
)

type AccountService interface {
	CreateAccount(account request.Account) (*domain.User, *response.Group, error)
}

type accountService struct {
	userService      UserService
	groupService     GroupService
	userGroupService UserGroupService
}

func NewAccountService(userService UserService, groupService GroupService, userGroupService UserGroupService) *accountService {
	return &accountService{
		userService:      userService,
		groupService:     groupService,
		userGroupService: userGroupService,
	}
}

func (s *accountService) CreateAccount(account request.Account) (*domain.User, *response.Group, error) {

	user := domain.User{
		UUID:     account.UUID,
		Name:     account.Name,
		Username: account.Username,
		Version:  1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser, err := s.userService.Create(ctx, user)
	if err != nil {
		slog.Error("Error when creating new user account", "error", err.Error())
		return nil, nil, err
	}

	newGroup, err := s.groupService.CreateDefaultGroup(ctx, *newUser)
	if err != nil {
		slog.Error("Error when creating default group to user", "error", err.Error())
		return nil, nil, err
	}

	userGroup := domain.UserGroup{
		UserId:  newUser.ID,
		GroupId: newGroup.ID,
	}
	s.userGroupService.Create(ctx, userGroup)

	return newUser, newGroup, nil
}
