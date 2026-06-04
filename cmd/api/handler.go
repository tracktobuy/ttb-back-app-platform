package main

import (
	"github.com/tracktobuy/ttb-back-app-platform/internal"
	"github.com/tracktobuy/ttb-back-app-platform/internal/handler"
)

type Handler struct {
	AccountHandler *handler.AccountHandler
	GroupHandler   *handler.GroupHandler
	ItemHandler    *handler.ItemHandler
	UserHandler    *handler.UserHandler
}

func CreateHandlers(services internal.Service) *Handler {
	return &Handler{
		AccountHandler: handler.NewAccountHandler(services),
		GroupHandler:   handler.NewGroupHandler(services),
		ItemHandler:    handler.NewItemHandler(services),
		UserHandler:    handler.NewUserHandler(services),
	}
}
