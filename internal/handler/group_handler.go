package handler

import (
	"context"
	"net/http"

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
	mux.HandleFunc("PUT /groups", h.Update)
}

func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {

	// 019e717e-7d9d-7529-9938-5e83bae19036
	groupId := helper.ReadParam(r, "groupId")
	group, err := h.service.GetByUUID(context.Background(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var updateRequest dto.UpdateGroupRequest

	err = helper.ReadJSON(w, r, &updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group.UUID = groupId
	group.Name = updateRequest.Name
	group.Budget = updateRequest.Budget
	group.BudgetCurrency = updateRequest.BudgetCurrency
	group.Version += 1

	updatedGroup, err := h.service.Update(context.Background(), groupId, *group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{"data": updatedGroup})
}
