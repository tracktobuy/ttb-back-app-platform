package internal

import (
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	UserRepo  repository.UserRepository
	GroupRepo repository.GroupRepository
	ItemRepo  repository.ItemRepository
	StoreRepo repository.StoreRepository
}

func CreateRepositories(db *mongo.Database) Repository {
	userRepo := repository.NewUserRepo(db)
	groupRepo := repository.NewGroupRepo(db)
	itemRepo := repository.NewItemRepo(db)
	storeRepo := repository.NewStoreRepo(db)

	return Repository{
		UserRepo:  userRepo,
		GroupRepo: groupRepo,
		ItemRepo:  itemRepo,
		StoreRepo: storeRepo,
	}
}
