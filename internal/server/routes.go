package server

import (
	"api.us4ever/internal/logger"
	"api.us4ever/internal/middleware"
	"api.us4ever/internal/server/routes"
	"github.com/gofiber/fiber/v2"
)

var (
	routesLogger *logger.Logger
)

func init() {
	var err error
	routesLogger, err = logger.New("routes")
	if err != nil {
		panic("failed to initialize routes logger: " + err.Error())
	}
}

func (s *FiberServer) RegisterFiberRoutes() {
	// 应用中间件
	s.App.Use(middleware.CORSMiddleware())
	s.App.Use(middleware.RequestIDMiddleware())
	s.App.Use(middleware.NewLoggingMiddleware())
	s.App.Use(middleware.NewPathBasedRateLimiter(100))
	s.App.Use(middleware.RecoveryMiddleware())

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

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) ErrorHandler(c *fiber.Ctx) error {
	panic("error")
}

func (s *FiberServer) AppConfigHandler(c *fiber.Ctx) error {
	return c.JSON(s.cfg)
}

func (s *FiberServer) UserListHandler(c *fiber.Ctx) error {
	users, err := s.DbClient.Client().User.Query().All(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(users)
}
