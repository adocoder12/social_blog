package handler

import (
	"log/slog"
)

type Application struct {
	logger *slog.Logger
}

func NewApplication(logger *slog.Logger) *Application {
	return &Application{logger: logger}
}
