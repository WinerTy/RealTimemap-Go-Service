package main

import (
	"log/slog"
	"os"
	"realtimemap-service/internal/app/repository/category"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/database/postgres"
)

const (
	envLocal = "local"
	envDev   = "Dev"
	envProd  = "Prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	storage, err := postgres.NewStorage(cfg.Database.BuildURL())
	if err != nil {
		log.Error("could not connect to database", err)
		os.Exit(1)
	}
	defer storage.Db.Close()

	// TODO Это на тест :)
	repo := category.NewPgCategoryRepository(storage.Db)

	items, err := repo.GetAll()
	if err != nil {
		log.Error("could not get items", err)
		os.Exit(1)
	}
	for _, item := range items {
		log.Info("Getting item", item)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return log
}
