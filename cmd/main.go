package main

import (
	"log"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	userHandler := handler.NewUserHandler()
	userHandler.Routes(mux)

	log.Printf("Server is running on port %s", ":8080")
	http.ListenAndServe(":8080", mux)
}
