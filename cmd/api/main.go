package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adocoder12/social_blog/internal/config"
	"github.com/adocoder12/social_blog/internal/db"
	"github.com/adocoder12/social_blog/internal/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to the database. Give up if it takes more than 5 seconds.
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(dbCtx, cfg.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close() // runs after the server has shut down
	logger.Info("database connection established")

	if err := db.Migrate(cfg.DatabaseConfig); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	logger.Info("database migrations verified")

	app := handler.NewApplication(logger)

	server := &http.Server{
		Addr:         ":" + cfg.SrvPort,
		Handler:      app.SetupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("API running", "port", cfg.SrvPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server failed: %w", err)
	case <-ctx.Done():
		logger.Info("shutting down")
	}

	// Restore default signal behavior: a second Ctrl+C now kills immediately.
	stop()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.Info("server stopped")
	return nil
}
