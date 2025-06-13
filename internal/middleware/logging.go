package middleware

import (
	"time"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingConfig defines the configuration for logging middleware
type LoggingConfig struct {
	// Logger instance to use
	Logger *logger.Logger

	// SkipPaths defines paths to skip logging
	SkipPaths []string

	// SkipSuccessfulRequests skips logging for successful requests (2xx status codes)
	SkipSuccessfulRequests bool

	// LogRequestBody enables logging of request body (be careful with sensitive data)
	LogRequestBody bool

	// LogResponseBody enables logging of response body (be careful with large responses)
	LogResponseBody bool

	// MaxBodySize limits the size of logged request/response bodies
	MaxBodySize int
}

// DefaultLoggingConfig returns a default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	httpLogger, err := logger.New("http")
	if err != nil {
		panic("failed to create http logger: " + err.Error())
	}

	return LoggingConfig{
		Logger:                 httpLogger,
		SkipPaths:              []string{"/health", "/metrics", "/favicon.ico"},
		SkipSuccessfulRequests: false,
		LogRequestBody:         true,
		LogResponseBody:        true,
		MaxBodySize:            1024, // 1KB
	}
}

// NewLoggingMiddleware creates a new logging middleware with the given configuration
func NewLoggingMiddleware(config ...LoggingConfig) fiber.Handler {
	cfg := DefaultLoggingConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Ensure logger is set
	if cfg.Logger == nil {
		httpLogger, err := logger.New("http")
		if err != nil {
			panic("failed to create http logger: " + err.Error())
		}
		cfg.Logger = httpLogger
	}

	return func(c *fiber.Ctx) error {
		// Skip logging for specified paths
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		// Generate request ID if not present
		requestID := GetRequestID(c)

		start := time.Now()

		// Prepare request fields
		requestFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		// Add query parameters if present
		if queryString := c.Context().QueryArgs().String(); queryString != "" {
			requestFields = append(requestFields, zap.String("query", queryString))
		}

		// Log request body if enabled
		if cfg.LogRequestBody && len(c.Body()) > 0 {
			body := c.Body()
			if len(body) > cfg.MaxBodySize {
				requestFields = append(requestFields, zap.String("request_body", string(body[:cfg.MaxBodySize])+"...[truncated]"))
			} else {
				requestFields = append(requestFields, zap.String("request_body", string(body)))
			}
		}

		// Log incoming request
		cfg.Logger.Info("incoming request", requestFields...)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Prepare response fields
		responseFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("duration", utils.SmartDurationFormat(duration)),
			zap.Int("size", len(c.Response().Body())),
		}

		// Add error information if present
		if err != nil {
			responseFields = append(responseFields, zap.Error(err))
		}

		// Log response body if enabled
		if cfg.LogResponseBody && len(c.Response().Body()) > 0 {
			body := c.Response().Body()
			if len(body) > cfg.MaxBodySize {
				responseFields = append(responseFields, zap.String("response_body", string(body[:cfg.MaxBodySize])+"...[truncated]"))
			} else {
				responseFields = append(responseFields, zap.String("response_body", string(body)))
			}
		}

		// Determine log level based on status code and configuration
		logLevel := getLogLevel(status, err, cfg.SkipSuccessfulRequests)

		// Log response
		switch logLevel {
		case "debug":
			cfg.Logger.Debug("request completed", responseFields...)
		case "info":
			cfg.Logger.Info("request completed", responseFields...)
		case "warn":
			cfg.Logger.Warn("request completed with warning", responseFields...)
		case "error":
			cfg.Logger.Error("request completed with error", responseFields...)
		}

		return err
	}
}

// getLogLevel determines the appropriate log level based on status code and error
func getLogLevel(status int, err error, skipSuccessful bool) string {
	if err != nil {
		return "error"
	}

	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warn"
	case status >= 200 && status < 300:
		if skipSuccessful {
			return "debug"
		}
		return "info"
	default:
		return "info"
	}
}

// RequestIDMiddleware adds a request ID to each request if not present
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Store in locals for easy access
		c.Locals("request_id", requestID)

		return c.Next()
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	// Fallback to header
	return c.Get("X-Request-ID")
}
