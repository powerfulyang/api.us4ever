package middleware

import (
	"context"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/errors"
	"api.us4ever/internal/logger"
	"api.us4ever/internal/utils"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
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
	healthLogger, err := logger.New("health")
	if err != nil {
		panic("failed to create health logger: " + err.Error())
	}

	return &HealthMiddleware{
		checkers: make(map[string]HealthChecker),
		logger:   healthLogger,
		timeout:  5 * time.Second,
	}
}

// AddChecker adds a health checker
func (h *HealthMiddleware) AddChecker(name string, checker HealthChecker) {
	if name == "" || checker == nil {
		h.logger.Warn("invalid health checker",
			zap.String("name", name),
			zap.Bool("checker", checker != nil),
		)
		return
	}

	h.checkers[name] = checker
	h.logger.Info("health checker added",
		zap.String("name", name),
	)
}

// SetTimeout sets the timeout for health checks
func (h *HealthMiddleware) SetTimeout(timeout time.Duration) {
	if timeout <= 0 {
		h.logger.Warn("invalid timeout, using default",
			zap.Duration("timeout", timeout),
			zap.Duration("default", h.timeout),
		)
		return
	}
	h.timeout = timeout
}

// Handler returns the health check handler
func (h *HealthMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
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

		response.Duration = utils.SmartDurationFormat(time.Since(startTime))

		// Set appropriate HTTP status code
		statusCode := fiber.StatusOK
		if !overallHealthy {
			statusCode = fiber.StatusServiceUnavailable
		}

		// Log health check result
		h.logger.Info("health check completed",
			zap.String("status", response.Status),
			zap.String("duration", response.Duration),
			zap.Int("components", len(response.Components)),
		)

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

			h.logger.Error("component health check failed",
				zap.String("component", name),
				zap.Error(err),
			)
		} else {
			h.logger.Debug("component health check passed",
				zap.String("component", name),
			)
		}
	case <-ctx.Done():
		status.Status = "unhealthy"
		status.Error = "health check timeout"

		h.logger.Error("component health check timeout",
			zap.String("component", name),
			zap.Duration("timeout", h.timeout),
		)
	}

	return status
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
	client *elasticsearch.Client
}

// NewElasticsearchHealthCheckerWithClient creates a new Elasticsearch health checker with client
func NewElasticsearchHealthCheckerWithClient(client *elasticsearch.Client) *ElasticsearchHealthChecker {
	return &ElasticsearchHealthChecker{client: client}
}

// Health checks Elasticsearch health
func (e *ElasticsearchHealthChecker) Health(ctx context.Context) error {
	if e.client == nil {
		return errors.NewInternalError("elasticsearch client is not initialized", nil)
	}

	// Try to cast to the expected ES client type and perform a ping
	// This is a simplified implementation - in practice you'd import the ES client type
	// For now, we'll assume the client is healthy if it's not nil
	// TODO: Implement actual ES ping when the client interface is properly defined
	// Check Elasticsearch connection
	pingResp, err := e.client.Ping(e.client.Ping.WithContext(ctx))
	if err != nil {
		return errors.NewInternalError("elasticsearch connection error", err)
	}
	defer pingResp.Body.Close()
	if pingResp.IsError() {
		return errors.NewInternalError("elasticsearch service unavailable", nil)
	}

	// If both checks pass
	return nil
}
