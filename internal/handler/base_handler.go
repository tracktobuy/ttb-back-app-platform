package handler

import (
	"github.com/tracktobuy/ttb-back-app-platform/internal/logger"
)

type BaseHandler interface {
	SetLogger(log logger.Logger)
}
