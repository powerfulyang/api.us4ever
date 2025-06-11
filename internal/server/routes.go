package server

import (
	"net/http"
	"time"

	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

// requestTimerMiddleware logs the time taken for each request.
func requestTimerMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	duration := time.Since(start)

	// Log details
	// Use c.Response().StatusCode() which is available after c.Next()
	routesLogger.Info("request completed", logger.Fields{
		"method":   c.Method(),
		"path":     c.Path(),
		"status":   c.Response().StatusCode(),
		"duration": duration.String(),
	})

	return err // Return the error reported by handlers
}

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Apply Request Timer middleware
	s.App.Use(requestTimerMiddleware)

	s.App.Get("/", s.HelloWorldHandler)
	internal := s.Group("/internal")

	internal.Get("/health", s.healthHandler)
	internal.Get("/app-config", s.AppConfigHandler)
	internal.Get("/user/list", s.UserListHandler)
	// Add the route for searching keeps
	internal.Get("/keeps/search", s.searchKeepsHandler)
	// Add the route for searching moments
	internal.Get("/moments/search", s.searchMomentsHandler)
	internal.Post("/keeps/reindex", s.reindexKeepsHandler)
	internal.Post("/moments/reindex", s.reindexMomentsHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	// Check database connection
	if err := s.DbClient.Health(c.Context()); err != nil {
		routesLogger.Error("health check failed: database connection error", logger.Fields{
			"error": err.Error(),
		})
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
		routesLogger.Error("health check failed: Elasticsearch ping error", logger.Fields{
			"error": err.Error(),
		})
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch connection error",
			"details": err.Error(),
		})
	}
	defer pingResp.Body.Close()
	if pingResp.IsError() {
		routesLogger.Error("health check failed: Elasticsearch ping returned error status", logger.Fields{
			"response": pingResp.String(),
		})
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
