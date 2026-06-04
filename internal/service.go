package internal

import (
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type Service struct {
	UserService    service.UserService
	GroupService   service.GroupService
	AccountService service.AccountService
	ItemService    service.ItemService
	StoreService   service.StoreService
}

func CreateServices(repo Repository) Service {

	userService := service.NewUserService(repo.UserRepo)
	groupService := service.NewGroupService(repo.GroupRepo)
	itemService := service.NewItemService(repo.ItemRepo)
	storeService := service.NewStoreService(repo.StoreRepo)

	return Service{
		UserService:    userService,
		GroupService:   groupService,
		AccountService: service.NewAccountService(userService, groupService),
		ItemService:    itemService,
		StoreService:   storeService,
	}
}
