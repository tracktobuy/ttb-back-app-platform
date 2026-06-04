package handler

import (
	"context"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemHandler struct {
	services internal.Service
}

func NewItemHandler(services internal.Service) *ItemHandler {

	return &ItemHandler{
		services: services,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {

	var request request.ItemRequest
	err := helper.ReadJSON(w, r, &request)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"data": nil, "message": "Invalid request", "error": err.Error()})
		return
	}

	group, err := h.services.GroupService.Get(context.Background(), request.GroupId)
	if err != nil {
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

	newStore, err := h.services.StoreService.Create(context.Background(), store)
	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	item := domain.Item{
		Title:  request.Title,
		Images: []string{request.Image},
		Groups: []primitive.ObjectID{group.ID},
		Stores: []primitive.ObjectID{newStore.ID},
	}

	newItem, err := h.services.ItemService.Create(context.Background(), item)

	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	response := h.formatItemResponse(newItem)
	response.Groups = []string{group.UUID}
	response.Stores = []string{newStore.UUID}

	helper.WriteJSON(w, http.StatusCreated, response)

}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	groupID := r.URL.Query().Get("groupId")

	group, err := h.services.GroupService.Get(context.Background(), groupID)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

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

	helper.WriteJSON(w, http.StatusOK, responseItems)
}

func (h *ItemHandler) formatItemResponse(item *domain.Item) *response.Item {
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
