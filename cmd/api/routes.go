package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Accounts
	mux.HandleFunc("POST /accounts", app.handlers.AccountHandler.CreateAccount)

	// Groups
	mux.HandleFunc("PUT /groups/{groupUUID}", app.handlers.GroupHandler.Update)
	mux.HandleFunc("GET /groups/{groupUUID}", app.handlers.GroupHandler.Get)
	mux.HandleFunc("GET /groups/{groupUUID}/labels", app.handlers.GroupHandler.GetLabels)
	mux.HandleFunc("GET /groups/{groupUUID}/items", app.handlers.GroupHandler.GetItems)

	// Items
	mux.HandleFunc("POST /items", app.handlers.ItemHandler.Create)
	mux.HandleFunc("GET /items", app.handlers.ItemHandler.GetAll)
	mux.HandleFunc("GET /items/{itemId}", app.handlers.ItemHandler.GetByUUID)

	// Stores
	mux.HandleFunc("GET /items/{itemId}/stores", app.handlers.StoreHandler.GetStoresByItemId)
	mux.HandleFunc("PUT /stores/{storeId}", app.handlers.StoreHandler.Update)
	mux.HandleFunc("DELETE /stores/{storeId}", app.handlers.StoreHandler.Delete)

	// Users
	mux.HandleFunc("POST /users", app.handlers.UserHandler.Create)

	return mux
}
