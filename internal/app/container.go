package app

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/config"
	"realtimemap-service/internal/database/postgres"
	categoryDomain "realtimemap-service/internal/domain/category"
	markDomain "realtimemap-service/internal/domain/mark"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"
	categoryRepo "realtimemap-service/internal/repository/category/postgres"
	markRepo "realtimemap-service/internal/repository/mark/postgres"
	"realtimemap-service/internal/service/category"
	markService "realtimemap-service/internal/service/mark"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Container Структура хранения зависимостей проекта, при появлении зависимостей добавить ее сюда и в метод NewContainer
type Container struct {
	Config *config.Config
	Logger *slog.Logger

	DbPool *pgxpool.Pool
	Redis  *redis.Client
	Cache  cache.Store

	CategoryRepository categoryDomain.Repository
	MarkRepository     markDomain.Repository

	CategoryService categoryDomain.Service
	MarkService     markDomain.Service
}

// NewContainer Фабрика, которая собирает все зависимости проекта в единый контейнер
func NewContainer(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Container, error) {
	pool, err := postgres.NewStorage(ctx, cfg.Database.BuildURL())
	if err != nil {
		logger.Error("could not connect to database", sl.Err(err))
		return nil, err
	}

	cacheStore, redisCli := config.InitCache(ctx, cfg)

	CategoryRepository := categoryRepo.NewPgCategoryRepository(pool)
	CategoryService := category.NewServiceCategory(CategoryRepository)

	MarkRepository := markRepo.NewPgMarkRepository(pool)
	MarkService := markService.NewService(MarkRepository)

	return &Container{
		Config: cfg,
		Logger: logger,

		DbPool: pool,
		Redis:  redisCli,
		Cache:  cacheStore,

		CategoryRepository: CategoryRepository,
		CategoryService:    CategoryService,

		MarkRepository: MarkRepository,
		MarkService:    MarkService,
	}, nil
}

func (c *Container) Close() error {
	if err := c.Redis.Close(); err != nil {
		c.Logger.Error("could not close redis", sl.Err(err))
		return err
	}
	return nil
}
