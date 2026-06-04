package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
)

type GroupHandler struct {
	service service.GroupService
}

func NewGroupHandler(groupService service.GroupService) *GroupHandler {
	return &GroupHandler{
		service: groupService,
	}
}

func (h *GroupHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("PUT /groups/{groupId}", h.Update)
	mux.HandleFunc("GET /groups/{groupId}", h.Get)
}

func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {

	groupId := r.PathValue("groupId")

	group, err := h.service.Get(context.Background(), groupId)

	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	helper.WriteJSON(w, http.StatusOK, envelope{"data": group})
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {

	groupId := r.PathValue("groupId")
	group, err := h.service.Get(context.Background(), groupId)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	var reqGroup request.Group

	err = helper.ReadJSON(w, r, &reqGroup)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"message": "Invalid request body", "data": nil})
		return
	}

	group.UUID = groupId
	group.Name = reqGroup.Name
	group.Budget = reqGroup.Budget
	group.BudgetCurrency = reqGroup.BudgetCurrency

	updatedGroup, err := h.service.Update(context.Background(), groupId, reqGroup)
	if err != nil {
		slog.Error("Error when updating group", "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"message": "Failed to update group", "data": envelope{"error": err.Error()}})
		return
	}

	helper.WriteJSON(w, http.StatusOK, envelope{"data": updatedGroup})
}
