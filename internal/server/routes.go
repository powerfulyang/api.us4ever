package server

import (
	"net/http"
	"time"

	"log"

	"strings"

	"api.us4ever/internal/task/image"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// requestTimerMiddleware logs the time taken for each request.
func requestTimerMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	duration := time.Since(start)

	// Log details
	// Use c.Response().StatusCode() which is available after c.Next()
	log.Printf("[%s] %s %d - %s", c.Method(), c.Path(), c.Response().StatusCode(), duration)

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
	internal.Get("/keeps/search", s.SearchKeepsHandler)
	// Add the route for searching moments
	internal.Get("/moments/search", s.SearchMomentsHandler)
	internal.Post("/keeps/reindex", s.ReindexKeepsHandler)
	internal.Post("/moments/reindex", s.ReindexMomentsHandler)
	// OCR
	internal.Get("/ocr/:id", s.OCRHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	// Check database connection
	if err := s.db.Health(c.Context()); err != nil {
		log.Printf("Health check failed: Database connection error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection error",
			"details": err.Error(),
		})
	}

	// Check Elasticsearch connection
	if s.esClient == nil {
		log.Printf("Health check failed: Elasticsearch client is not initialized.")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch client not initialized",
		})
	}
	pingResp, err := s.esClient.Ping(s.esClient.Ping.WithContext(c.Context()))
	if err != nil {
		log.Printf("Health check failed: Elasticsearch ping error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Elasticsearch connection error",
			"details": err.Error(),
		})
	}
	defer pingResp.Body.Close()
	if pingResp.IsError() {
		log.Printf("Health check failed: Elasticsearch ping returned error status: %s", pingResp.String())
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
	users, err := s.db.Client().User.Query().All(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(users)
}

func (s *FiberServer) OCRHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image ID is required",
		})
	}

	err := image.ProcessSingleImageOCR(c.Context(), s.db, id)
	if err != nil {
		// Check for specific error types and return appropriate status codes
		if strings.Contains(err.Error(), "image not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// For other errors, return 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "OCR processing completed successfully",
		"id":      id,
	})
}
