package middleware

import (
	"context"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/errors"
	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v2"
)

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	Health(ctx context.Context) error
}

// HealthStatus represents the health status of a component
type HealthStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status     string         `json:"status"`
	Timestamp  time.Time      `json:"timestamp"`
	Components []HealthStatus `json:"components"`
	Duration   string         `json:"duration"`
}

// HealthMiddleware provides health check functionality
type HealthMiddleware struct {
	checkers map[string]HealthChecker
	logger   *logger.Logger
	timeout  time.Duration
}

// NewHealthMiddleware creates a new health middleware
func NewHealthMiddleware() *HealthMiddleware {
	return &HealthMiddleware{
		checkers: make(map[string]HealthChecker),
		logger:   logger.New("health"),
		timeout:  5 * time.Second,
	}
}

// AddChecker adds a health checker
func (h *HealthMiddleware) AddChecker(name string, checker HealthChecker) {
	if name == "" || checker == nil {
		h.logger.Warn("invalid health checker", logger.Fields{
			"name":    name,
			"checker": checker != nil,
		})
		return
	}

	h.checkers[name] = checker
	h.logger.Info("health checker added", logger.Fields{
		"name": name,
	})
}

// SetTimeout sets the timeout for health checks
func (h *HealthMiddleware) SetTimeout(timeout time.Duration) {
	if timeout <= 0 {
		h.logger.Warn("invalid timeout, using default", logger.Fields{
			"timeout": timeout,
			"default": h.timeout,
		})
		return
	}
	h.timeout = timeout
}

// Handler returns the health check handler
func (h *HealthMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Context(), h.timeout)
		defer cancel()

		response := HealthResponse{
			Timestamp:  startTime,
			Components: make([]HealthStatus, 0, len(h.checkers)),
		}

		overallHealthy := true

		// Check all registered health checkers
		for name, checker := range h.checkers {
			status := h.checkComponent(ctx, name, checker)
			response.Components = append(response.Components, status)

			if status.Status != "healthy" {
				overallHealthy = false
			}
		}

		// Set overall status
		if overallHealthy {
			response.Status = "healthy"
		} else {
			response.Status = "unhealthy"
		}

		response.Duration = time.Since(startTime).String()

		// Set appropriate HTTP status code
		statusCode := fiber.StatusOK
		if !overallHealthy {
			statusCode = fiber.StatusServiceUnavailable
		}

		// Log health check result
		h.logger.Info("health check completed", logger.Fields{
			"status":     response.Status,
			"duration":   response.Duration,
			"components": len(response.Components),
		})

		return c.Status(statusCode).JSON(response)
	}
}

// checkComponent checks the health of a single component
func (h *HealthMiddleware) checkComponent(ctx context.Context, name string, checker HealthChecker) HealthStatus {
	status := HealthStatus{
		Name:   name,
		Status: "healthy",
	}

	// Create a channel to receive the result
	done := make(chan error, 1)

	// Run health check in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- errors.NewInternalError("health check panicked", nil)
			}
		}()
		done <- checker.Health(ctx)
	}()

	// Wait for result or timeout
	select {
	case err := <-done:
		if err != nil {
			status.Status = "unhealthy"
			status.Error = err.Error()

			h.logger.Error("component health check failed", logger.Fields{
				"component": name,
				"error":     err.Error(),
			})
		} else {
			h.logger.Debug("component health check passed", logger.Fields{
				"component": name,
			})
		}
	case <-ctx.Done():
		status.Status = "unhealthy"
		status.Error = "health check timeout"

		h.logger.Error("component health check timeout", logger.Fields{
			"component": name,
			"timeout":   h.timeout,
		})
	}

	return status
}

// ReadinessHandler returns a simple readiness check handler
func (h *HealthMiddleware) ReadinessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	}
}

// LivenessHandler returns a simple liveness check handler
func (h *HealthMiddleware) LivenessHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "alive",
			"timestamp": time.Now(),
		})
	}
}

// DatabaseHealthChecker implements HealthChecker for database
type DatabaseHealthChecker struct {
	db database.Service
}

// NewDatabaseHealthChecker creates a new database health checker
func NewDatabaseHealthChecker(db database.Service) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

// Health checks database health
func (d *DatabaseHealthChecker) Health(ctx context.Context) error {
	if d.db == nil {
		return errors.NewDatabaseError("database service is nil", nil)
	}

	return d.db.Health(ctx)
}

// ElasticsearchHealthChecker implements HealthChecker for Elasticsearch
type ElasticsearchHealthChecker struct {
	// Add ES client when needed
}

// NewElasticsearchHealthChecker creates a new Elasticsearch health checker
func NewElasticsearchHealthChecker() *ElasticsearchHealthChecker {
	return &ElasticsearchHealthChecker{}
}

// Health checks Elasticsearch health
func (e *ElasticsearchHealthChecker) Health(ctx context.Context) error {
	// TODO: Implement ES health check
	return nil
}
