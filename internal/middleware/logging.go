package middleware

import (
	"strconv"
	"time"

	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		LogRequestBody:         false,
		LogResponseBody:        false,
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
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Store request ID in context for other middleware/handlers
		c.Locals("request_id", requestID)

		start := time.Now()

		// Prepare request fields
		requestFields := logger.Fields{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       path,
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		}

		// Add query parameters if present
		if queryString := c.Context().QueryArgs().String(); queryString != "" {
			requestFields["query"] = queryString
		}

		// Log request body if enabled
		if cfg.LogRequestBody && len(c.Body()) > 0 {
			body := c.Body()
			if len(body) > cfg.MaxBodySize {
				requestFields["request_body"] = string(body[:cfg.MaxBodySize]) + "...[truncated]"
			} else {
				requestFields["request_body"] = string(body)
			}
		}

		// Log incoming request
		cfg.Logger.Info("incoming request", requestFields)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Prepare response fields
		responseFields := logger.Fields{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       path,
			"status":     status,
			"duration":   duration.String(),
			"size":       len(c.Response().Body()),
		}

		// Add error information if present
		if err != nil {
			responseFields["error"] = err.Error()
		}

		// Log response body if enabled
		if cfg.LogResponseBody && len(c.Response().Body()) > 0 {
			body := c.Response().Body()
			if len(body) > cfg.MaxBodySize {
				responseFields["response_body"] = string(body[:cfg.MaxBodySize]) + "...[truncated]"
			} else {
				responseFields["response_body"] = string(body)
			}
		}

		// Determine log level based on status code and configuration
		logLevel := getLogLevel(status, err, cfg.SkipSuccessfulRequests)

		// Log response
		switch logLevel {
		case "debug":
			cfg.Logger.Debug("request completed", responseFields)
		case "info":
			cfg.Logger.Info("request completed", responseFields)
		case "warn":
			cfg.Logger.Warn("request completed with warning", responseFields)
		case "error":
			cfg.Logger.Error("request completed with error", responseFields)
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

// CorrelationIDMiddleware adds correlation ID support for distributed tracing
func CorrelationIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		correlationID := c.Get("X-Correlation-ID")
		if correlationID == "" {
			// Use request ID as correlation ID if not provided
			if requestID := c.Get("X-Request-ID"); requestID != "" {
				correlationID = requestID
			} else {
				correlationID = uuid.New().String()
			}
			c.Set("X-Correlation-ID", correlationID)
		}

		// Store in locals for easy access
		c.Locals("correlation_id", correlationID)

		return c.Next()
	}
}

// MetricsMiddleware provides basic request metrics
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Log metrics (in a real application, you might send these to a metrics system)
		logger.Info("request metrics", logger.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   status,
			"duration": duration.Milliseconds(),
			"size":     len(c.Response().Body()),
		})

		return err
	}
}

// SecurityHeadersMiddleware adds common security headers
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Add cache control for API responses
		if c.Path() != "/health" && c.Path() != "/metrics" {
			c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")
		}

		return c.Next()
	}
}

// RateLimitInfo represents rate limit information
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// AddRateLimitHeaders adds rate limit headers to the response
func AddRateLimitHeaders(c *fiber.Ctx, info RateLimitInfo) {
	c.Set("X-RateLimit-Limit", strconv.Itoa(info.Limit))
	c.Set("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))
	c.Set("X-RateLimit-Reset", strconv.FormatInt(info.Reset.Unix(), 10))
}
