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
	logger    *slog.Logger
	config    *config.Config
	validator *validator.Validate
	version   string

	posts    *services.PostService

}

func NewApplication(
	logger *slog.Logger,
	cfg *config.Config,
	version string,
	posts *services.PostService,
) *Application {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Report errors using json tag names (user_id) instead of Go field names (UserID).
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Application{
		logger:    logger,
		config:    cfg,
		validator: v,
		version:   version,
		posts:     posts,
	}
}