package app

import (
	"lesson-manager-server/internal/config"
	"lesson-manager-server/internal/storage"
	"lesson-manager-server/internal/storage/postgres"
	"log/slog"
	"os"
)

func SetupLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger
	switch cfg.Env {
	case "dev":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func SetupStorage(cfg *config.Database) (*storage.Storage, error) {
	var s *storage.Storage
	var err error
	switch cfg.Driver {
	case "postgres":
		s, err = postgres.New(cfg)
		//case "sqlite":
		//	sqlite.New(...)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}
