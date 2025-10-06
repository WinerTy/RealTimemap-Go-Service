package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"realtimemap-service/internal/app"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/pkg/logger/sl"
	v1 "realtimemap-service/internal/transport/http/v1"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	ctx := context.Background()

	container, err := app.NewContainer(ctx, cfg, log)
	if err != nil {
		log.Error("Error creating container", sl.Err(err))
		os.Exit(1)
	}
	defer container.DbPool.Close()
	defer container.Redis.Close()

	r := gin.Default()
	v1.InitV1Routers(r, container)

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
