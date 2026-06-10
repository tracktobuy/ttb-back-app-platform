package handler

import (
	"context"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type GroupHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type groupHandler struct {
	log     logger.Logger
	service internal.Service
}

func NewGroupHandler(service internal.Service) GroupHandler {

	log := logger.NewLogger()
	log.SetHandlerName("GroupHandler")

	return &groupHandler{
		service: service,
		log:     log,
	}
}
func (h *groupHandler) Get(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Get")

	groupId := r.PathValue("groupId")

	group, err := h.service.GroupService.Get(context.Background(), groupId)

	h.log.Info("get group by id", "groupId", groupId)

	if err != nil {
		h.log.Error("get group by id", "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	h.log.Info("get group by id success", "groupId", groupId, "response", group)

	helper.WriteJSON(w, http.StatusOK, envelope{"data": group})
}

func (h *groupHandler) Update(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Update")

	groupId := r.PathValue("groupId")
	h.log.Info("update group", "groupId", groupId)

	group, err := h.service.GroupService.Get(context.Background(), groupId)

	if err != nil {
		h.log.Error("update group failed", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	var reqGroup request.Group

	err = helper.ReadJSON(w, r, &reqGroup)
	if err != nil {
		h.log.Error("update group failed", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"message": "Invalid request body", "data": nil})
		return
	}

	group.UUID = groupId
	group.Name = reqGroup.Name
	group.Budget = reqGroup.Budget
	group.BudgetCurrency = reqGroup.BudgetCurrency

	h.log.Info("update group request", "groupId", groupId, "request", group)

	updatedGroup, err := h.service.GroupService.Update(context.Background(), groupId, reqGroup)
	if err != nil {
		h.log.Error("update group failed", "groupId", groupId, "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"message": "Failed to update group", "data": envelope{"error": err.Error()}})
		return
	}

	h.log.Info("update group request success", "groupId", groupId, "request", updatedGroup)
	helper.WriteJSON(w, http.StatusOK, envelope{"data": updatedGroup})
}
