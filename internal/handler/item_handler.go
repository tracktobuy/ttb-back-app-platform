package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type ItemHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	log      logger.Logger
	services internal.Service
}

func NewItemHandler(services internal.Service) ItemHandler {

	log := logger.NewLogger()
	log.SetHandlerName("ItemHandler")

	return &itemHandler{
		services: services,
		log:      log,
	}
}

func (h *itemHandler) Create(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Create")

	var request request.ItemRequest
	err := helper.ReadJSON(w, r, &request)

	h.log.Info("create new item request", "request", request)

	if err != nil {
		h.log.Error("create item failed", "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"data": nil, "message": "Invalid request", "error": err.Error()})
		return
	}

	group, err := h.services.GroupService.Get(context.Background(), request.GroupId)
	if err != nil {
		h.log.Error("finding group with id failed", "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

	ac := helper.GetCookie(w, r)

	user, err := h.services.UserService.GetById(context.Background(), ac.UserObjectID())
	if err != nil {
		h.log.Error("user not found", "error", err.Error())
	}

	item := domain.Item{
		Title:     request.Title,
		Images:    []string{request.Image},
		CreatedBy: user.ID,
	}

	newItem, err := h.services.ItemService.Create(context.Background(), item)

	if err != nil {
		h.log.Error("creating item failed", "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	store := domain.Store{
		UUID:       helper.GenerateUUIDV7(),
		Price:      request.Price,
		Currency:   request.Currency,
		Domain:     request.Domain,
		Name:       request.Store,
		BestOption: false,
		Item:       newItem.ID,
		CreatedBy:  user.ID,
		URL:        request.URL,
	}

	_, err = h.services.StoreService.Create(context.Background(), store)
	if err != nil {
		h.log.Error("creating store for item failed", "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	groupItem := domain.GroupItem{
		GroupId: group.ID,
		ItemId:  newItem.ID,
	}

	_, err = h.services.GroupItemService.Create(context.Background(), groupItem)
	if err != nil {
		h.log.Error("creating junction with group and item", "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	h.log.Info("creating new item success", "response", newItem)
	helper.WriteJSON(w, http.StatusCreated, envelope{"data": newItem})

}

func (h *itemHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("GetAll")

	acc := helper.GetCookie(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	h.log.Info("get all items by user uuid", "userUUID", acc.UserUUID)

	items, err := h.services.ItemService.GetAllByUserId(ctx, acc.UserObjectID())
	if err != nil {
		h.log.Error("error when get items for user", "userUUID", acc.UserUUID, "error", err.Error())
		helper.InternalServerError(w, err)
	}

	helper.WriteJSON(w, http.StatusOK, envelope{"data": items})

}
