package handler

import (
	"context"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	log     logger.Logger
	service internal.Service
}

func NewUserHandler(service internal.Service) UserHandler {
	log := logger.NewLogger()
	log.SetHandlerName("UserHandler")

	return &userHandler{service: service, log: log}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Create")

	var userRequest request.Account

	err := helper.ReadJSON(w, r, &userRequest)

	h.log.Info("create new user", "request", userRequest)

	if err != nil {
		helper.BadRequest(w, r, err)
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

	h.log.Info("create new user success", "response", userRequest)

	helper.WriteJSON(w, http.StatusCreated, envelope{"data": user})
}
