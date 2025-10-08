package v1

import (
	"realtimemap-service/internal/app"
	"realtimemap-service/internal/transport/http/v1/mark"

	"realtimemap-service/internal/transport/http/v1/category"

	"github.com/gin-gonic/gin"
)

// InitV1Routers Функция регистрации v1 GIN роутеров в основное приложение.
// Принимет Приложение GIN + Контейнер зависимостей
func InitV1Routers(g *gin.Engine, container *app.Container) {
	r := g.Group("/api/v1")
	category.InitCategoryRoutes(r.Group("/category"), container.CategoryService, container.Cache)
	mark.InitMarkRoutes(r.Group("/marks"), container.MarkService)
}
