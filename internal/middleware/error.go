package middleware

import (
	"fmt"

	"api.us4ever/internal/errors"
	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   ErrorDetail `json:"error"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// ErrorHandlerConfig defines the configuration for error handling
type ErrorHandlerConfig struct {
	// Logger instance to use
	Logger *logger.Logger

	// IncludeStackTrace includes stack trace in development mode
	IncludeStackTrace bool

	// CustomErrorMap maps specific errors to custom responses
	CustomErrorMap map[error]ErrorResponse

	// DefaultErrorMessage is used when error message should be hidden
	DefaultErrorMessage string
}

// DefaultErrorHandlerConfig returns a default error handler configuration
func DefaultErrorHandlerConfig() ErrorHandlerConfig {
	errorLogger, err := logger.New("error")
	if err != nil {
		panic("failed to create error logger: " + err.Error())
	}

	return ErrorHandlerConfig{
		Logger:              errorLogger,
		IncludeStackTrace:   false,
		CustomErrorMap:      make(map[error]ErrorResponse),
		DefaultErrorMessage: "An internal error occurred",
	}
}

// NewErrorHandler creates a new error handling middleware
func NewErrorHandler(config ...ErrorHandlerConfig) fiber.ErrorHandler {
	cfg := DefaultErrorHandlerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Ensure logger is set
	if cfg.Logger == nil {
		errorLogger, err := logger.New("error")
		if err != nil {
			panic("failed to create error logger: " + err.Error())
		}
		cfg.Logger = errorLogger
	}

	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}

		// Get request ID for tracing
		requestID := getRequestID(c)

		// Log the error with context
		logError(cfg.Logger, err, c, requestID)

		// Determine response based on error type
		response := buildErrorResponse(err, cfg, requestID)
		statusCode := determineStatusCode(err)

		// Set content type
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		return c.Status(statusCode).JSON(response)
	}
}

// getRequestID extracts request ID from context
func getRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	// Fallback to header
	return c.Get("X-Request-ID")
}

// logError logs the error with appropriate context
func logError(log *logger.Logger, err error, c *fiber.Ctx, requestID string) {
	fields := logger.Fields{
		"error":      err.Error(),
		"method":     c.Method(),
		"path":       c.Path(),
		"ip":         c.IP(),
		"user_agent": c.Get("User-Agent"),
	}

	if requestID != "" {
		fields["request_id"] = requestID
	}

	// Add query parameters if present
	if queryString := c.Context().QueryArgs().String(); queryString != "" {
		fields["query"] = queryString
	}

	// Check if it's an application error
	if appErr := errors.GetAppError(err); appErr != nil {
		fields["error_type"] = appErr.Type
		if appErr.Code > 0 {
			fields["error_code"] = appErr.Code
		}

		// Log at appropriate level based on error type
		switch appErr.Type {
		case "ValidationError", "NotFoundError":
			log.Warn("client error occurred", fields)
		case "DatabaseError", "ElasticsearchError":
			log.Error("external service error occurred", fields)
		default:
			log.Error("application error occurred", fields)
		}
	} else {
		// Log as error for unknown error types
		log.Error("unhandled error occurred", fields)
	}
}

// buildErrorResponse constructs the error response
func buildErrorResponse(err error, cfg ErrorHandlerConfig, requestID string) ErrorResponse {
	// Check for custom error mapping first
	if customResponse, exists := cfg.CustomErrorMap[err]; exists {
		customResponse.TraceID = requestID
		return customResponse
	}

	response := ErrorResponse{
		TraceID: requestID,
	}

	// Handle application errors
	if appErr := errors.GetAppError(err); appErr != nil {
		response.Error = ErrorDetail{
			Type:    appErr.Type,
			Message: appErr.Message,
			Code:    appErr.Code,
		}
		return response
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		response.Error = ErrorDetail{
			Type:    "HTTPError",
			Message: fiberErr.Message,
			Code:    fiberErr.Code,
		}
		return response
	}

	// Handle generic errors
	response.Error = ErrorDetail{
		Type:    "InternalError",
		Message: cfg.DefaultErrorMessage,
		Code:    500,
	}

	return response
}

// determineStatusCode determines the HTTP status code for the error
func determineStatusCode(err error) int {
	// Handle application errors
	if appErr := errors.GetAppError(err); appErr != nil {
		if appErr.Code > 0 {
			return appErr.Code
		}

		// Default status codes based on error type
		switch appErr.Type {
		case "ValidationError":
			return fiber.StatusBadRequest
		case "NotFoundError":
			return fiber.StatusNotFound
		case "ConfigError":
			return fiber.StatusInternalServerError
		case "DatabaseError", "ElasticsearchError":
			return fiber.StatusServiceUnavailable
		default:
			return fiber.StatusInternalServerError
		}
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		return fiberErr.Code
	}

	// Default to internal server error
	return fiber.StatusInternalServerError
}

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() fiber.Handler {
	recoveryLogger, err := logger.New("recovery")
	if err != nil {
		panic("failed to create recovery logger: " + err.Error())
	}

	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				requestID := getRequestID(c)

				// Log the panic
				recoveryLogger.Error("panic recovered", logger.Fields{
					"panic":      fmt.Sprintf("%v", r),
					"request_id": requestID,
					"method":     c.Method(),
					"path":       c.Path(),
					"ip":         c.IP(),
				})

				// Create error response
				response := ErrorResponse{
					Error: ErrorDetail{
						Type:    "InternalError",
						Message: "An internal error occurred",
						Code:    500,
					},
					TraceID: requestID,
				}

				// Send error response
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
				c.Status(fiber.StatusInternalServerError).JSON(response)
			}
		}()

		return c.Next()
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := getRequestID(c)

		response := ErrorResponse{
			Error: ErrorDetail{
				Type:    "NotFoundError",
				Message: fmt.Sprintf("Route '%s %s' not found", c.Method(), c.Path()),
				Code:    404,
			},
			TraceID: requestID,
		}

		return c.Status(fiber.StatusNotFound).JSON(response)
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := getRequestID(c)

		response := ErrorResponse{
			Error: ErrorDetail{
				Type:    "MethodNotAllowedError",
				Message: fmt.Sprintf("Method '%s' not allowed for route '%s'", c.Method(), c.Path()),
				Code:    405,
			},
			TraceID: requestID,
		}

		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}
