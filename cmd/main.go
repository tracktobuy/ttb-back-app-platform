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

	repository := repository.NewUserRepo(db)
	service := service.NewUserService(repository)
	userHandler := handler.NewUserHandler(service)
	userHandler.Routes(mux)

	log.Printf("Server is running on port %s", ":8080")
	http.ListenAndServe(":8080", mux)
}
