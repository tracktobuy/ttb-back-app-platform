package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
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

	groupUUID := r.PathValue("groupUUID")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	group, err := h.service.GroupService.Get(ctx, groupUUID)

	h.log.Info("get group by id", "groupUUID", groupUUID)

	if err != nil {
		h.log.Error("get group by id", "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	user, err := h.service.UserService.GetById(ctx, group.CreatedBy.ID)
	if err != nil {
		h.log.Error("get group by id", "error", err.Error())
		details := response.ClientErrorDetails{
			Field: "groupUUID",
			Issue: "Group not found with groupUUID" + groupUUID,
		}
		helper.NotFound(w, r, err, details)
		return
	}

	group.CreatedBy = response.User{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Name:     user.Name,
	}

	h.log.Info("get group by id success", "groupUUID", groupUUID, "response", group)

	helper.WriteJSON(w, http.StatusOK, envelope{"data": group})
}

func (h *groupHandler) Update(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Update")

	groupUUID := r.PathValue("groupUUID")
	h.log.Info("update group", "groupUUID", groupUUID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	group, err := h.service.GroupService.Get(ctx, groupUUID)

	if err != nil {
		h.log.Error("update group failed", "groupUUID", groupUUID, "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"message": "Group not found", "data": nil})
		return
	}

	user, err := h.service.UserService.GetById(ctx, group.CreatedBy.ID)
	if err != nil {
		h.log.Error("update group failed", "groupUUID", groupUUID, "error", err.Error())
		details := response.ClientErrorDetails{
			Field: "groupUUID",
			Issue: "Group not found with groupUUID" + groupUUID,
		}
		helper.NotFound(w, r, err, details)
		return
	}

	var reqGroup request.Group

	err = helper.ReadJSON(w, r, &reqGroup)
	if err != nil {
		h.log.Error("update group failed", "groupUUID", groupUUID, "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"message": "Invalid request body", "data": nil})
		return
	}

	group.UUID = groupUUID
	group.Name = reqGroup.Name
	group.Budget = reqGroup.Budget
	group.BudgetCurrency = reqGroup.BudgetCurrency

	h.log.Info("update group request", "groupUUID", groupUUID, "request", group)

	updatedGroup, err := h.service.GroupService.Update(ctx, groupUUID, reqGroup)
	if err != nil {
		h.log.Error("update group failed", "groupUUID", groupUUID, "error", err.Error())
		helper.WriteJSON(w, http.StatusInternalServerError, envelope{"message": "Failed to update group", "data": envelope{"error": err.Error()}})
		return
	}

	updatedGroup.CreatedBy = response.User{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Name:     user.Name,
	}

	h.log.Info("update group request success", "groupUUID", groupUUID, "request", updatedGroup)
	helper.WriteJSON(w, http.StatusOK, envelope{"data": updatedGroup})
}

func (h *groupHandler) GetLabels(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("GetLabels")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	groupUUID := r.PathValue("groupUUID")

	h.log.Info("listing labels for group", "groupUUID", groupUUID)

	data, err := h.service.GroupService.GetLabels(ctx, groupUUID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.log.Error("group not found", "groupUUID", groupUUID)
			details := response.ClientErrorDetails{
				Field: "groupUUID",
				Issue: "Group not found with groupUUID" + groupUUID,
			}
			helper.NotFound(w, r, err, details)
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

	groupUUID := r.PathValue("groupUUID")

	h.log.Info("listing items for group", "groupUUID", groupUUID)

	data, err := h.service.GroupService.GetItems(ctx, groupUUID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.log.Error("group not found", "groupUUID", groupUUID)
			details := response.ClientErrorDetails{
				Field: "groupUUID",
				Issue: "Group not found with groupUUID" + groupUUID,
			}
			helper.NotFound(w, r, err, details)
			return
		}

		h.log.Error("error when get items", "groupUUID", groupUUID)
		helper.InternalServerError(w, err)
		return
	}

	h.log.Info("listing items for group success", "groupUUID", groupUUID, "total items", len(data.Items))
	helper.WriteJSON(w, http.StatusOK, envelope{"data": data})
}
