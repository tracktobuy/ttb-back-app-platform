package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemHandler interface {
	BaseHandler
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type itemHandler struct {
	logger   *slog.Logger
	services internal.Service
}

func NewItemHandler(services internal.Service) ItemHandler {

	return &itemHandler{
		services: services,
	}
}

func (h *itemHandler) SetLogger(logger *slog.Logger) {
	h.logger = logger
}

func (h *itemHandler) Create(w http.ResponseWriter, r *http.Request) {

	var request request.ItemRequest
	err := helper.ReadJSON(w, r, &request)
	if err != nil {
		h.logger.Error("create item failed", "handler", "ItemHandler", "method", "Create", "error", err.Error())
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"data": nil, "message": "Invalid request", "error": err.Error()})
		return
	}

	h.logger.Info("finding group with id", "handler", "ItemHandler", "method", "Create", "groupId", request.GroupId)

	group, err := h.services.GroupService.Get(context.Background(), request.GroupId)
	if err != nil {
		h.logger.Error("finding group with id failed", "handler", "ItemHandler", "method", "Create", "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

	store := domain.Store{
		UUID:       helper.GenerateUUIDV7(),
		Price:      request.Price,
		Currency:   request.Currency,
		Domain:     request.Domain,
		Name:       request.Store,
		BestOption: false,
	}

	h.logger.Info("creating store for item", "handler", "ItemHandler", "method", "Create", "store", store)
	newStore, err := h.services.StoreService.Create(context.Background(), store)
	if err != nil {
		h.logger.Error("creating store for item failed", "handler", "ItemHandler", "method", "Create", "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	item := domain.Item{
		Title:  request.Title,
		Images: []string{request.Image},
		Groups: []primitive.ObjectID{group.ID},
		Stores: []primitive.ObjectID{newStore.ID},
	}

	h.logger.Info("creating item", "handler", "ItemHandler", "method", "Create", "item", store)
	newItem, err := h.services.ItemService.Create(context.Background(), item)

	if err != nil {
		h.logger.Error("creating item failed", "handler", "ItemHandler", "method", "Create", "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	response := h.formatItemResponse(newItem)
	response.Groups = []string{group.UUID}
	response.Stores = []string{newStore.UUID}

	h.logger.Info("creating item success", "handler", "ItemHandler", "method", "Create", "response", response)
	helper.WriteJSON(w, http.StatusCreated, envelope{"data": response})

}

func (h *itemHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	groupID := r.URL.Query().Get("groupId")

	h.logger.Info("find group by id", "handler", "ItemHandler", "method", "GetAll", "groupId", groupID)

	group, err := h.services.GroupService.Get(context.Background(), groupID)
	if err != nil {
		h.logger.Error("find group by id failed", "handler", "ItemHandler", "method", "GetAll", "groupId", groupID, "error", err.Error())
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

	h.logger.Info("find all items by group id", "handler", "ItemHandler", "method", "GetAll", "groupId", group.UUID)

	items, err := h.services.ItemService.GetAllByGroupID(context.Background(), group.ID)

	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	var responseItems []response.Item

	for _, item := range items {
		responseItem := h.formatItemResponse(&item)
		responseItems = append(responseItems, *responseItem)
	}

	if len(responseItems) == 0 {
		helper.WriteJSON(w, http.StatusOK, envelope{"data": []response.Item{}})
		return
	}

	helper.WriteJSON(w, http.StatusOK, envelope{"data": responseItems})
}

func (h *itemHandler) formatItemResponse(item *domain.Item) *response.Item {
	return &response.Item{
		UUID:      item.UUID,
		Title:     item.Title,
		Images:    item.Images,
		CreatedAt: helper.DateTime(item.CreatedAt),
		UpdatedAt: helper.DateTime(item.UpdatedAt),
		Labels:    item.Labels,
		Stores:    []string{},
		Groups:    []string{},
	}
}
