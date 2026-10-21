package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/adocoder12/social_blog/internal/config"
	"github.com/adocoder12/social_blog/internal/db"
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
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// Connect to the database. Give up if it takes more than 5 seconds.
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(dbCtx, cfg.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()
	logger.Info("database connection established")

	// Seeding hashes a password and runs ~800 inserts, so it gets its own, longer budget.
	seedCtx, cancelSeed := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelSeed()

	if err := db.Seed(seedCtx, pool); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}
	return nil
}
