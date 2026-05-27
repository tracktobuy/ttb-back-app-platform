package service

import (
	"context"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
)

type accountService struct {
	userService  *userService
	groupService *groupService
	ctx          context.Context
}

func NewAccountService(ctx context.Context, userService *userService, groupService *groupService) *accountService {
	return &accountService{
		userService:  userService,
		groupService: groupService,
		ctx:          ctx,
	}
}

func (s *accountService) CreateAccount(user domain.User) error {

	_, err := s.userService.Create(s.ctx, user)
	if err != nil {
		return err
	}

	return nil

}
