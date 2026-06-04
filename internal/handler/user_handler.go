package handler

import (
	"context"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
)

type UserHandler struct {
	service internal.Service
}

func NewUserHandler(service internal.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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
