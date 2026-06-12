package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Accounts
	mux.HandleFunc("POST /accounts", app.handlers.AccountHandler.CreateAccount)

	// Groups
	mux.HandleFunc("PUT /groups/{groupId}", app.handlers.GroupHandler.Update)
	mux.HandleFunc("GET /groups/{groupId}", app.handlers.GroupHandler.Get)

	// Items
	mux.HandleFunc("POST /items", app.handlers.ItemHandler.Create)
	mux.HandleFunc("GET /items", app.handlers.ItemHandler.GetAll)
	mux.HandleFunc("GET /items/{itemId}", app.handlers.ItemHandler.GetByUUID)

	// Stores
	mux.HandleFunc("GET /items/{itemId}/stores", app.handlers.StoreHandler.GetStoresByItemId)

	// Users
	mux.HandleFunc("POST /users", app.handlers.UserHandler.Create)

	return mux
}
