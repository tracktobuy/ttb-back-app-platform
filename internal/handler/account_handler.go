package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type accountHandler struct {
	userService  service.UserService
	groupService service.GroupService
}

func NewAccountHandler(ctx context.Context,
	userService service.UserService,
	groupService service.GroupService) *accountHandler {

	return &accountHandler{
		userService:  userService,
		groupService: groupService,
	}
}

func (h *accountHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /accounts", h.CreateAccount)
}

func (h *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	var userRequest request.Account
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, group, err := service.NewAccountService(ctx, h.userService, h.groupService).CreateAccount(newUser)

	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	response := NewAccountResponse(*user, group)

	helper.WriteJSON(w, http.StatusCreated, map[string]any{"data": response})

}

type userResponse struct {
	UUID     string           `json:"uuid"`
	Username string           `json:"username"`
	Name     string           `json:"name"`
	Groups   []response.Group `json:"groups"`
}

type newAccountResponse struct {
	User userResponse `json:"user"`
}

func NewAccountResponse(user domain.User, group *response.Group) newAccountResponse {

	groups := make([]response.Group, 0)
	if group != nil {
		groups = append(groups, *group)
	}

	return newAccountResponse{
		User: userResponse{
			UUID:     user.UUID,
			Username: user.Username,
			Name:     user.Name,
			Groups:   groups,
		},
	}
}
