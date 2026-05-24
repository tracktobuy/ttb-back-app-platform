package main

import (
	"log"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
	"github.com/tracktobuy/ttb-back-app-platform/internal/repository"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

func main() {
	mux := http.NewServeMux()

	repository := repository.NewUserRepo(nil)
	service := service.NewUserService(repository)
	userHandler := handler.NewUserHandler(service)
	userHandler.Routes(mux)

	log.Printf("Server is running on port %s", ":8080")
	http.ListenAndServe(":8080", mux)
}
