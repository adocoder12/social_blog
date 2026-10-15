package handler

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/adocoder12/social_blog/internal/config"
	"github.com/adocoder12/social_blog/internal/services"
	"github.com/go-playground/validator/v10"
)

type Application struct {
	logger *slog.Logger
	posts  *services.PostService
	config  *config.Config
	validator *validator.Validate
	version string
}

func NewApplication(logger *slog.Logger, posts *services.PostService ,cfg *config.Config, version string) *Application {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Application{
		logger:    logger,
		posts:     posts,
		config:    cfg,
		validator: v,
		version:   version,
	}
}
