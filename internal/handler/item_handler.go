package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/domain"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/request"
	"github.com/tracktobuy/ttb-back-app-platform/internal/dto/response"
	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemHandler struct {
	itemService  service.ItemServiceInterface
	groupService service.CrudService[domain.Group]
}

func NewItemHandler(itemService service.ItemServiceInterface,
	groupService service.CrudService[domain.Group]) *ItemHandler {

	return &ItemHandler{
		itemService:  itemService,
		groupService: groupService,
	}
}

func (h *ItemHandler) Routes(mux *http.ServeMux) {

	mux.HandleFunc("POST /items", h.Create)
	mux.HandleFunc("GET /items", h.GetAll)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {

	var request request.ItemRequest
	err := helper.ReadJSON(w, r, &request)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, envelope{"data": nil, "message": "Invalid request", "error": err.Error()})
		return
	}

	group, err := h.groupService.Get(context.Background(), request.GroupId)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

	item := domain.Item{
		Title:  request.Title,
		Images: []string{request.Image},
		Groups: []primitive.ObjectID{group.ID},
	}

	newItem, err := h.itemService.Create(context.Background(), item)

	if err != nil {
		helper.InternalServerError(w, err)
		return
	}

	response := h.formatItemResponse(newItem)
	response.Groups = []string{group.UUID}

	helper.WriteJSON(w, http.StatusCreated, response)

}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	groupID := r.URL.Query().Get("groupId")
	slog.Info("GetAll", "info", groupID)

	group, err := h.groupService.Get(context.Background(), groupID)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, envelope{"data": nil, "message": "Group not found", "error": err.Error()})
		return
	}

	items, err := h.itemService.GetAllByGroupID(context.Background(), group.ID)

	slog.Info("Length", "info", len(items))

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
