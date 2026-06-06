package handler

import (
	"log/slog"
)

type BaseHandler interface {
	SetLogger(logger *slog.Logger)
}
