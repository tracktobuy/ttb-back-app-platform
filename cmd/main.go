package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/config"
	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.LoadConfig()

	client := config.MongoConnect(ctx, cfg)
	defer client.Disconnect(context.Background())

	mux := http.NewServeMux()

	db := client.Database(cfg.MongoDB)

	userRepository := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	userHandler.Routes(mux)

	groupRepository := repository.NewGroupRepo(db)
	groupService := service.NewGroupServiceImplementation(groupRepository)

	accountHandler := handler.NewAccountHandler(ctx, userService, groupService)
	accountHandler.RegisterRoutes(mux)

	log.Printf("Server is running on port %s", ":8080")
	http.ListenAndServe(":8080", mux)
}
