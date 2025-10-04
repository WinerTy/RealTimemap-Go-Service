package main

import (
	"context"
	"log/slog"
	"os"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/database/postgres"
	"realtimemap-service/internal/pkg/logger/sl"
	repository "realtimemap-service/internal/repository/category/postgres"
	service "realtimemap-service/internal/service/category"
	http "realtimemap-service/internal/transport/http/v1/category" // TODO временно
	"time"

	"github.com/gin-contrib/cache/persistence"
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

	store := persistence.NewInMemoryStore(time.Minute) // TODO Мб сделаь самописный с помощью Redis?

	pool, err := postgres.NewStorage(ctx, cfg.Database.BuildURL())

	if err != nil {
		log.Error("could not connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.NewPgCategoryRepository(pool)
	serv := service.NewServiceCategory(repo)

	r := gin.Default()
	http.InitCategoryRoutes(r.Group("/api"), serv, store)

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
