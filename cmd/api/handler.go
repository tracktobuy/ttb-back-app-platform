package main

import (
	"log/slog"

	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	AccountHandler handler.AccountHandler
	GroupHandler   handler.GroupHandler
	ItemHandler    handler.ItemHandler
	UserHandler    handler.UserHandler
}

func CreateHandlers(db *mongo.Database) *Handler {

	services := internal.CreateServices(internal.CreateRepositories(db))

	return &Handler{
		AccountHandler: handler.NewAccountHandler(services),
		GroupHandler:   handler.NewGroupHandler(services),
		ItemHandler:    handler.NewItemHandler(services),
		UserHandler:    handler.NewUserHandler(services),
	}
}

func (h *Handler) WithLogger(logger *slog.Logger) *Handler {
	h.AccountHandler.SetLogger(logger)
	h.GroupHandler.SetLogger(logger)
	h.ItemHandler.SetLogger(logger)
	h.UserHandler.SetLogger(logger)
	return h
}
