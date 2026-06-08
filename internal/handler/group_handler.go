package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
)

type GroupHandler interface {
	BaseHandler
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type groupHandler struct {
	logger  *slog.Logger
	service internal.Service
}

func NewGroupHandler(service internal.Service) GroupHandler {
	return &groupHandler{
		service: service,
	}
}

func (h *groupHandler) SetLogger(logger *slog.Logger) {
	h.logger = logger
}

func (h *groupHandler) Get(w http.ResponseWriter, r *http.Request) {

	groupId := r.PathValue("groupId")

	group, err := h.service.GroupService.Get(context.Background(), groupId)

	h.logger.Info("get group by id", "handler", "GroupHandler", "method", "Get", "groupId", groupId)

	if err != nil {
		h.logger.Error("get group by id", "handler", "GroupHandler", "method", "Get", "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	h.logger.Info("get group by id success", "handler", "GroupHandler", "method", "Get", "groupId", groupId, "response", group)

	helper.WriteJSON(w, http.StatusOK, envelope{"data": group})
}

func (h *groupHandler) Update(w http.ResponseWriter, r *http.Request) {

	groupId := r.PathValue("groupId")
	h.logger.Info("update group", "handler", "GroupHandler", "method", "Update", "groupId", groupId)

	group, err := h.service.GroupService.Get(context.Background(), groupId)

	if err != nil {
		h.logger.Error("update group failed", "handler", "GroupHandler", "method", "Update", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	var reqGroup request.Group

	err = helper.ReadJSON(w, r, &reqGroup)
	if err != nil {
		h.logger.Error("update group failed", "handler", "GroupHandler", "method", "Update", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"message": "Invalid request body", "data": nil})
		return
	}

	group.UUID = groupId
	group.Name = reqGroup.Name
	group.Budget = reqGroup.Budget
	group.BudgetCurrency = reqGroup.BudgetCurrency

	h.logger.Info("update group request", "handler", "GroupHandler", "method", "Update", "groupId", groupId, "request", group)

	updatedGroup, err := h.service.GroupService.Update(context.Background(), groupId, reqGroup)
	if err != nil {
		h.logger.Error("update group failed", "handler", "GroupHandler", "method", "Update", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"message": "Failed to update group", "data": envelope{"error": err.Error()}})
		return
	}

	h.logger.Info("update group request success", "handler", "GroupHandler", "method", "Update", "groupId", groupId, "request", updatedGroup)
	helper.WriteJSON(w, http.StatusOK, envelope{"data": updatedGroup})
}
