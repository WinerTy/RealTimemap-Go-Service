package main

import (
	"context"
	"log/slog"
	"os"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/database/postgres"
	v1 "realtimemap-service/internal/handler/http/v1"
	repository "realtimemap-service/internal/repository/category/postgres"
	service "realtimemap-service/internal/service/category"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "Dev"
	envProd  = "Prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	ctx := context.Background()
	pool, err := postgres.NewStorage(ctx, cfg.Database.BuildURL())
	if err != nil {
		log.Error("could not connect to database", err)
		os.Exit(1)
	}
	defer pool.Close()

	// TODO Это на тест :)
	repo := repository.NewPgCategoryRepository(pool)
	serv := service.NewServiceCategory(repo)
	r := gin.Default()
	v1.InitCategoryRoutes(r.Group("/api"), serv)
	r.Run()
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
