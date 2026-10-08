package handler

import (
	"log/slog"

	"github.com/adocoder12/social_blog/internal/store"
)

type Application struct {
	logger *slog.Logger
	store  store.Storage
}

func NewApplication(logger *slog.Logger, store store.Storage) *Application {
	return &Application{logger: logger, store: store}
}
