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

type GroupHandler struct {
	service service.GroupServiceInterface
}

func NewGroupHandler(groupService service.GroupServiceInterface) *GroupHandler {
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

	groupResponse := h.formatGroupResponse(group)
	helper.WriteJSON(w, http.StatusOK, envelope{"data": groupResponse})
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {

	groupId := r.PathValue("groupId")
	group, err := h.service.Get(context.Background(), groupId)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	var updateRequest dto.UpdateGroupRequest

	err = helper.ReadJSON(w, r, &updateRequest)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"message": "Invalid request body", "data": nil})
		return
	}

	group.UUID = groupId
	group.Name = updateRequest.Name
	group.Budget = updateRequest.Budget
	group.BudgetCurrency = updateRequest.BudgetCurrency

	updatedGroup, err := h.service.Update(context.Background(), groupId, *group)
	if err != nil {
		slog.Error("Error when updating group", "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"message": "Failed to update group", "data": envelope{"error": err.Error()}})
		return
	}

	response := h.formatGroupResponse(updatedGroup)
	helper.WriteJSON(w, http.StatusOK, envelope{"data": response})
}

func (h *GroupHandler) formatGroupResponse(group *domain.Group) dto.GroupResponse {
	return dto.GroupResponse{
		UUID:           group.UUID,
		Name:           group.Name,
		Budget:         group.Budget,
		BudgetCurrency: group.BudgetCurrency,
		CreatedAt:      helper.DateTime(group.CreatedAt),
		UpdatedAt:      helper.DateTime(group.UpdatedAt),
	}
}
