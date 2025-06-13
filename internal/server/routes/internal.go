package routes

import (
	"api.us4ever/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type InternalRoutes struct {
	app      *fiber.App
	cfg      interface{}
	dbClient interface{}
	esClient interface{}
}

func NewInternalRoutes(app *fiber.App, cfg interface{}, dbClient, esClient interface{}) *InternalRoutes {
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
	internal.Get("/app-config", r.AppConfigHandler)
	internal.Get("/user/list", r.UserListHandler)
}

func (r *InternalRoutes) AppConfigHandler(c *fiber.Ctx) error {
	return c.JSON(r.cfg)
}

func (r *InternalRoutes) UserListHandler(c *fiber.Ctx) error {
	users, err := r.dbClient.(interface{ Client() interface{} }).Client().User.Query().All(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(users)
}
