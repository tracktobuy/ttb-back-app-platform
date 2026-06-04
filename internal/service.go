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
	ItemService    service.ItemService
}

func CreateServices(userRepo repository.CrudRepository[domain.User],
	groupRepo repository.GroupRepository,
	itemRepo repository.ItemRepository) Service {

	userService := service.NewUserService(userRepo)
	groupService := service.NewGroupService(groupRepo)
	itemService := service.NewItemService(itemRepo)

	return Service{
		UserService:    userService,
		GroupService:   groupService,
		AccountService: service.NewAccountService(userService, groupService),
		ItemService:    itemService,
	}
}
