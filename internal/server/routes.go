package server

import (
	"time"

	"api.us4ever/internal/middleware"
	"api.us4ever/internal/server/routes"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// 应用中间件
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:       3600,
	}))
	s.App.Use(requestid.New())
	s.App.Use(middleware.NewLoggingMiddleware())
	s.App.Use(middleware.RecoveryMiddleware())
	s.App.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// 注册基础路由
	routes.RegisterBaseRoutes(s.App)

	// 注册内部路由
	internalRoutes := routes.NewInternalRoutes(s.App, s.cfg, s.DbClient, s.EsClient)
	internalRoutes.Register()

	// 注册搜索路由
	searchRoutes := routes.NewSearchRoutes(s.App)
	searchRoutes.Register()

	// 注册重索引路由
	reindexRoutes := routes.NewReindexRoutes(s.App)
	reindexRoutes.Register()
}
