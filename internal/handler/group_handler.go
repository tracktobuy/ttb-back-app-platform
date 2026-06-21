package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)

	GetLabels(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
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

func (h *groupHandler) GetLabels(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("GetLabels")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	groupUUID := r.PathValue("groupId")

	h.log.Info("listing labels for group", "groupUUID", groupUUID)

	data, err := h.service.GroupService.GetLabels(ctx, groupUUID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.log.Error("group not found", "groupUUID", groupUUID)
			helper.NotFound(w, err)
			return
		}

		h.log.Error("error when get labels", "groupUUID", groupUUID)
		helper.InternalServerError(w, err)
		return
	}

	h.log.Info("listing labels for group success", "groupUUID", groupUUID, "total labels", len(data.Labels))

	helper.WriteJSON(w, http.StatusOK, envelope{"data": data})

}

func (h *groupHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	h.log.SetMethodName("GetItems")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	groupUUID := r.PathValue("groupId")

	h.log.Info("listing items for group", "groupUUID", groupUUID)

	data, err := h.service.GroupService.GetItems(ctx, groupUUID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.log.Error("group not found", "groupUUID", groupUUID)
			helper.NotFound(w, err)
			return
		}

		h.log.Error("error when get items", "groupUUID", groupUUID)
		helper.InternalServerError(w, err)
		return
	}

	h.log.Info("listing items for group success", "groupUUID", groupUUID, "total items", len(data.Items))
	helper.WriteJSON(w, http.StatusOK, envelope{"data": data})
}
