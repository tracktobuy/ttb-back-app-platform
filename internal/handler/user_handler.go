package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
)

type UserHandler interface {
	BaseHandler
	Create(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	logger  *slog.Logger
	service internal.Service
}

func NewUserHandler(service internal.Service) UserHandler {
	return &userHandler{service: service}
}

func (h *userHandler) SetLogger(logger *slog.Logger) {
	h.logger = logger
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userRequest request.Account

	err := helper.ReadJSON(w, r, &userRequest)

	if err != nil {
		helper.BadRequest(w, err)
		return
	}

	item := domain.User{
		UUID:     userRequest.UUID,
		Name:     userRequest.Name,
		Username: userRequest.Username,
		Version:  1,
	}

	user, err := h.service.UserService.Create(context.Background(), item)

	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, envelope{"data": user})
}
