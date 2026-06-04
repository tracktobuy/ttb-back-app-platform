package main

import (
	"log"
	"net/http"
)

func (app *application) run() {

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Printf("Server is running on port %s", ":8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
