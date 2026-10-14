package handler

import (
	"log/slog"

	"github.com/adocoder12/social_blog/internal/config"
	"github.com/adocoder12/social_blog/internal/services"
)

type Application struct {
	logger *slog.Logger
	posts  *services.PostService
	config  *config.Config
	version string
}

func NewApplication(logger *slog.Logger, posts *services.PostService ,cfg *config.Config, version string) *Application {
	return &Application{logger: logger, posts: posts, config: cfg , version: version}
}
