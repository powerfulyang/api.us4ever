package server

import (
	"net/http"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	// Apply middleware from middleware package
	s.App.Use(middleware.CORSMiddleware())
	s.App.Use(middleware.RequestIDMiddleware())
	s.App.Use(middleware.NewLoggingMiddleware())
	s.App.Use(middleware.NewPathBasedRateLimiter(1))

	// Apply error handling middleware
	s.App.Use(middleware.RecoveryMiddleware())

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/err", s.ErrorHandler)
	internal := s.Group("/internal")

	// Health endpoint without rate limiting
	internal.Get("/health", s.healthHandler)
	internal.Get("/app-config", s.AppConfigHandler)
	internal.Get("/user/list", s.UserListHandler)

	// Search endpoints with enhanced validation and rate limiting
	searchGroup := internal.Group("/search")

	// Search routes with new paths
	searchGroup.Get("/keeps", s.searchKeepsHandler)
	searchGroup.Get("/moments", s.searchMomentsHandler)

	// Keep old routes for backward compatibility
	internal.Get("/keeps/search", s.searchKeepsHandler)
	internal.Get("/moments/search", s.searchMomentsHandler)

	// Reindex endpoints
	internal.Post("/keeps/reindex", s.reindexKeepsHandler)
	internal.Post("/moments/reindex", s.reindexMomentsHandler)
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

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	// Check database connection
	if err := s.DbClient.Health(c.Context()); err != nil {
		routesLogger.Error("health check failed: database connection error",
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection error",
			"details": err.Error(),
		})
	}

	// Check Elasticsearch connection
	if s.EsClient == nil {
		routesLogger.Error("health check failed: Elasticsearch client is not initialized")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch client not initialized",
		})
	}
	pingResp, err := s.EsClient.Ping(s.EsClient.Ping.WithContext(c.Context()))
	if err != nil {
		routesLogger.Error("health check failed: Elasticsearch ping error",
			zap.Error(err),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch connection error",
			"details": err.Error(),
		})
	}
	defer pingResp.Body.Close()
	if pingResp.IsError() {
		routesLogger.Error("health check failed: Elasticsearch ping returned error status",
			zap.String("response", pingResp.String()),
		)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch service unavailable",
			"details": pingResp.String(), // Include ES response string for details
		})
	}

	// If both checks pass
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Health check passed (Database & Elasticsearch OK)",
	})
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
