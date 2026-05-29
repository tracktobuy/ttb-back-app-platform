package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type accountHandler struct {
	ctx          context.Context
	userService  service.CrudService[domain.User]
	groupService service.GroupServiceInterface
}

func NewAccountHandler(ctx context.Context,
	userService service.CrudService[domain.User],
	groupService service.GroupServiceInterface) *accountHandler {

	return &accountHandler{
		ctx:          ctx,
		userService:  userService,
		groupService: groupService,
	}
}

func (h *accountHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /accounts", h.CreateAccount)
}

func (h *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var userRequest dto.NewUserRequest
	err := helper.ReadJSON(w, r, &userRequest)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	newUser := domain.User{
		UUID:     userRequest.UUID,
		Name:     userRequest.Name,
		Username: userRequest.Username,
		Version:  1,
	}

	slog.Info("Requesting user creation")
	user, group, err := service.NewAccountService(h.ctx, h.userService, h.groupService).CreateAccount(newUser)

	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	response := dto.NewAccountResponse(*user, *group)

	helper.WriteJSON(w, http.StatusCreated, map[string]any{"data": response})

}
