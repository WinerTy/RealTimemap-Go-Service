package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/database/postgres"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"
	repository "realtimemap-service/internal/repository/category/postgres"
	service "realtimemap-service/internal/service/category"
	myhttp "realtimemap-service/internal/transport/http/v1/category" // TODO временно
	"syscall"
	"time"

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
	var store cache.Store
	redisClient, err := config.SetupRedis(ctx, cfg)
	if err != nil {
		log.Error("Error setting up redis client")
		store = cache.NewNoOpCache()
	} else {
		store = cache.NewRedisCache(redisClient)

	}

	pool, err := postgres.NewStorage(ctx, cfg.Database.BuildURL())
	if err != nil {
		log.Error("could not connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.NewPgCategoryRepository(pool)
	serv := service.NewServiceCategory(repo)

	r := gin.Default()
	myhttp.InitCategoryRoutes(r.Group("/api"), serv, store)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server", sl.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", sl.Err(err))
	}
	slog.Info("Server exiting")
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
