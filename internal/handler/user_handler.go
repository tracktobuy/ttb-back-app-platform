package handler

import (
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type UserHandler struct {
	service service.CrudService[domain.User]
}

func NewUserHandler(service service.CrudService[domain.User]) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.Create)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusCreated, map[string]any{"data": "somethings"})
}
