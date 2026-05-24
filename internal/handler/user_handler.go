package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type UserHandler struct {
	service service.CrudService[domain.User]
}

type userRequest struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewUserHandler(service service.CrudService[domain.User]) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.Create)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	var userRequest userRequest
	err := helper.ReadJSON(w, r, &userRequest)

	if err != nil {
		helper.BadRequest(w, err)
		return
	}

	item := domain.User{
		UUID:     userRequest.UUID,
		Name:     userRequest.Name,
		Username: userRequest.Username,
	}

	user, err := h.service.Create(context.Background(), item)

	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	log.Printf("Parsed request body: %+v", userRequest)
	log.Printf("Created user: %+v", user)

	helper.WriteJSON(w, http.StatusCreated, map[string]any{"data": user})
}
