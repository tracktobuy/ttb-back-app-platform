package service

import (
	"context"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type AccountService interface {
	CreateAccount(account request.Account) (*domain.User, *response.Group, error)
}

type accountService struct {
	userService      UserService
	groupService     GroupService
	userGroupService UserGroupService
	log              logger.Logger
}

func NewAccountService(userService UserService, groupService GroupService, userGroupService UserGroupService) *accountService {

	log := logger.NewLogger()
	log.SetServiceName("AccountService")

	return &accountService{
		userService:      userService,
		groupService:     groupService,
		userGroupService: userGroupService,
		log:              log,
	}
}

func (s *accountService) CreateAccount(account request.Account) (*domain.User, *response.Group, error) {

	s.log.SetMethodName("CreateAccount")

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
		s.log.Error("Error when creating new user account", "error", err.Error())
		return nil, nil, err
	}

	newGroup, err := s.groupService.CreateDefaultGroup(ctx, *newUser)
	if err != nil {
		s.log.Error("Error when creating default group to user", "error", err.Error())
		return nil, nil, err
	}

	userGroup := domain.UserGroup{
		UserId:  newUser.ID,
		GroupId: newGroup.ID,
	}
	s.userGroupService.Create(ctx, userGroup)

	return newUser, newGroup, nil
}
