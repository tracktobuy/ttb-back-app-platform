package internal

import (
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type Service struct {
	UserService    service.UserService
	GroupService   service.GroupService
	AccountService service.AccountService
}

func NewService(userRepo repository.CrudRepository[domain.User],
	groupRepo repository.GroupRepository) Service {

	userService := service.NewUserService(userRepo)
	groupService := service.NewGroupService(groupRepo)

	return Service{
		UserService:    userService,
		GroupService:   groupService,
		AccountService: service.NewAccountService(userService, groupService),
	}
}
