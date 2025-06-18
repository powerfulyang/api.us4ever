package middleware

import (
	"fmt"
	"time"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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

	return func(c fiber.Ctx) error {
		// Skip logging for specified paths
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		// Generate request ID if not present
		requestID := GetRequestID(c)

		// Record start time
		start := time.Now()

		// Prepare request fields
		requestFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.String("ip", GetRealIP(c)),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		// Add query parameters if present
		if allQueries := c.Queries(); len(allQueries) > 0 {
			for key, value := range allQueries {
				queryString := fmt.Sprintf("key=%s, value=%s", key, value)
				requestFields = append(requestFields, zap.String("query", queryString))
			}
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

		// Store the original handler to be executed after our middleware
		err := c.Next()

		// Calculate duration after request processing
		duration := time.Since(start)

		// 重要：在所有处理（包括错误处理）完成后获取最终的状态码
		// 这将确保获取到的是经过错误处理后的真实状态码
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

		// Determine log level based on final status code and configuration
		logLevel := getLogLevel(status, err, cfg.SkipSuccessfulRequests)

		// Log response with appropriate level based on the final status code
		switch logLevel {
		case "debug":
			cfg.Logger.Debug("request completed", responseFields...)
		case "info":
			cfg.Logger.Info("request completed", responseFields...)
		case "warn":
			cfg.Logger.Warn("request completed with warning", responseFields...)
		case "error":
			cfg.Logger.Warn("request completed with error", responseFields...)
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

// GetRequestID extracts request ID from context
func GetRequestID(c fiber.Ctx) string {
	return requestid.FromContext(c)
}
