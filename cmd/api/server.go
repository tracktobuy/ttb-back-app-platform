package main

import (
	"log"
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/config"
)

func (app *application) run(cfg *config.Config) {

	srv := &http.Server{
		Addr:    ":" + cfg.ApiServerPort,
		Handler: app.routes(),
	}

	log.Printf("Server is running on port :%s", cfg.ApiServerPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
