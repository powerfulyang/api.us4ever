package routes

import (
	"api.us4ever/internal/config"
	"api.us4ever/internal/database"
	"api.us4ever/internal/middleware"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v3"
)

type InternalRoutes struct {
	app      *fiber.App
	cfg      *config.AppConfig
	dbClient database.Service
	esClient *elasticsearch.Client
}

func NewInternalRoutes(app *fiber.App, cfg *config.AppConfig, dbClient database.Service, esClient *elasticsearch.Client) *InternalRoutes {
	return &InternalRoutes{
		app:      app,
		cfg:      cfg,
		dbClient: dbClient,
		esClient: esClient,
	}
}

func (r *InternalRoutes) Register() {
	internal := r.app.Group("/internal")

	// 初始化健康检查中间件
	healthMiddleware := middleware.NewHealthMiddleware()

	// 添加数据库健康检查
	if r.dbClient != nil {
		dbHealthChecker := middleware.NewDatabaseHealthChecker(r.dbClient)
		healthMiddleware.AddChecker("database", dbHealthChecker)
	}

	// 添加 Elasticsearch 健康检查
	if r.esClient != nil {
		esHealthChecker := middleware.NewElasticsearchHealthCheckerWithClient(r.esClient)
		healthMiddleware.AddChecker("elasticsearch", esHealthChecker)
	}

	// 健康检查端点
	internal.Get("/health", healthMiddleware.Handler())
	internal.Get("/app-config", func(c fiber.Ctx) error {
		return c.JSON(r.cfg)
	})
}
