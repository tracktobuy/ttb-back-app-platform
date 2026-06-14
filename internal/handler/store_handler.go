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

type StoreHandler interface {
	GetStoresByItemId(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type storeHandler struct {
	log      logger.Logger
	services internal.Service
}

func NewStoreHandler(svc internal.Service) StoreHandler {

	log := logger.NewLogger()
	log.SetHandlerName("StoreHandler")

	return &storeHandler{
		services: svc,
		log:      log,
	}
}

func (h *storeHandler) GetStoresByItemId(w http.ResponseWriter, r *http.Request) {
	h.log.SetMethodName("GetStoresByItemId")

	itemUUID := r.PathValue("itemId")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	item, err := h.services.ItemService.GetByUUID(ctx, itemUUID)
	if err != nil {
		h.log.Error("finding item by uuid", "uuid", itemUUID, "error", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.NotFound(w, err)
			return
		}

		helper.InternalServerError(w, err)
	}

	h.log.Info("searching stores by item UUID", "uuid", itemUUID)
	stores, err := h.services.StoreService.GetStoresByItemId(ctx, item.ID)
	if err != nil {
		h.log.Error("something went wrong when listing stores", "itemUUID", itemUUID, "error", err.Error())
		helper.InternalServerError(w, err)
	}

	h.log.Info("stores found", "total stores", len(stores))

	helper.WriteJSON(w, http.StatusOK, envelope{"data": stores})

}

func (h *storeHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.log.SetMethodName("Update")

	storeUUID := r.PathValue("storeId")

	var storeReq request.Store
	err := helper.ReadJSON(w, r, &storeReq)
	if err != nil {
		h.log.Error("error when parsing JSON for store update", "error", err.Error())
		helper.BadRequest(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	resp, err := h.services.StoreService.Update(ctx, storeUUID, storeReq)
	if err != nil {
		h.log.Error("error when updating store", "storeUUID", storeUUID, "error", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.NotFound(w, err)
		}
		helper.InternalServerError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, envelope{"data": resp})
}

func (h *storeHandler) Delete(w http.ResponseWriter, r *http.Request) {

	h.log.SetMethodName("Delete")

	storeUUID := r.PathValue("storeId")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	store, err := h.services.StoreService.GetStoreByUUID(ctx, storeUUID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.log.Error("deleting store, store not found", "storeUUID", storeUUID, "error", err.Error())
			helper.NotFound(w, err)
			return
		}

		h.log.Error("an error occurred when deleting store", "storeUUID", storeUUID, "error", err.Error())
	}

	err = h.services.StoreService.Delete(ctx, store.ID)
	if err != nil {
		h.log.Error("error deleting store", "storeUUID", storeUUID, "error", err.Error())
		helper.InternalServerError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusNoContent, envelope{})
}
