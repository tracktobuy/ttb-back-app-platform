package handler

import (
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type AccountHandler interface {
	BaseHandler
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type accountHandler struct {
	log      logger.Logger
	services internal.Service
}

func NewAccountHandler(services internal.Service) AccountHandler {
	return &accountHandler{
		services: services,
	}
}

func (h *accountHandler) SetLogger(log logger.Logger) {
	h.log = log
}

func (h *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	h.log.SetHandlerName("AccountHandler")
	h.log.SetMethodName("CreateAccount")

	var accountRequest request.Account
	err := helper.ReadJSON(w, r, &accountRequest)

	h.log.Info("create account request", "payload", accountRequest)

	if err != nil {
		h.log.Error("create account failed", "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"error": err.Error()})
		return
	}

	user, group, err := h.services.AccountService.CreateAccount(accountRequest)

	if err != nil {
		h.log.Error("created account failed", "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"error": err.Error()})
		return
	}

	response := h.newAccountResponse(*user, group)
	h.log.Info("create account success", "response", response)
	helper.WriteJSON(w, http.StatusCreated, envelope{"data": response})

}

func (h *accountHandler) newAccountResponse(user domain.User, group *response.Group) response.NewAccountResponse {

	groups := make([]response.Group, 0)
	if group != nil {
		groups = append(groups, *group)
	}

	return response.NewAccountResponse{
		User: response.User{
			UUID:     user.UUID,
			Username: user.Username,
			Name:     user.Name,
			Groups:   groups,
		},
	}
}
