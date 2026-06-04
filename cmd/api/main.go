package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/config"
	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.LoadConfig()

	client := config.MongoConnect(ctx, cfg)
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.MongoDB)

	repositories := internal.CreateRepositories(db)
	services := internal.CreateServices(repositories)

	mux := http.NewServeMux()
	// Handlers
	userHandler := handler.NewUserHandler(services)
	userHandler.Routes(mux)

	groupHandler := handler.NewGroupHandler(services)
	groupHandler.Routes(mux)

	itemHandler := handler.NewItemHandler(services)
	itemHandler.Routes(mux)

	accountHandler := handler.NewAccountHandler(services)
	accountHandler.RegisterRoutes(mux)

	log.Printf("Server is running on port %s", ":8080")
	http.ListenAndServe(":8080", mux)
}
