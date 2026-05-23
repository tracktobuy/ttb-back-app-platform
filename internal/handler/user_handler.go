package handler

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type UserHandler struct {
	service service.CrudService[domain.User]
}

type userRequest struct {
	Hey string `json:"hey"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{service: nil}
}

func (h *UserHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", h.Create)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request: %s %s", r.Method, r.URL.Path)

	var userRequest userRequest
	helper.ReadJSON(r, &userRequest)

	log.Printf("Parsed request body: %+v", userRequest)

	helper.WriteJSON(w, http.StatusCreated, map[string]any{"data": "somethings"})
}
