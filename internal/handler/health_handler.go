package handler

import (
	"net/http"

	"github.com/tracktobuy/ttb-back-app-platform/internal/helper"
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type HealthHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
	log logger.Logger
}

func NewHealthHandler() HealthHandler {
	log := logger.NewLogger()
	log.SetHandlerName("HealthHandler")

	return &healthHandler{log: log}
}

func (h *healthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.log.SetMethodName("HealthCheck")

	err := helper.WriteJSON(w, http.StatusOK, envelope{"status": "UP"})
	if err != nil {
		h.log.Error("health check response failed", "error", err.Error())
	}
}
