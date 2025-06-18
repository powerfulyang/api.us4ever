package middleware

import (
	sErrors "errors"
	"fmt"

	"api.us4ever/internal/errors"
	"api.us4ever/internal/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
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

	return func(c fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}

		// Get request ID for tracing
		requestID := GetRequestID(c)

		// Determine response based on error type
		response := buildErrorResponse(err, cfg, requestID)
		statusCode := determineStatusCode(err)

		// Set content type
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		return c.Status(statusCode).JSON(response)
	}
}

// NewErrorMiddleware 创建一个错误处理中间件，用于在路由中注册
// 与 NewErrorHandler 不同的是，这个函数返回的是一个中间件处理函数
func NewErrorMiddleware(config ...ErrorHandlerConfig) fiber.Handler {
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

	return func(c fiber.Ctx) error {
		// 先调用下一个处理函数
		err := c.Next()

		// 如果有错误，则处理它
		if err != nil {
			// Get request ID for tracing
			requestID := GetRequestID(c)

			// Determine response based on error type
			response := buildErrorResponse(err, cfg, requestID)
			statusCode := determineStatusCode(err)

			// Set content type
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			// 返回JSON响应并设置状态码
			return c.Status(statusCode).JSON(response)
		}

		// 没有错误，返回nil
		return nil
	}
}

// buildErrorResponse constructs the error response
func buildErrorResponse(err error, cfg ErrorHandlerConfig, requestID string) ErrorResponse {
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
	var fiberErr *fiber.Error
	if sErrors.As(err, &fiberErr) {
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
	var fiberErr *fiber.Error
	if sErrors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	// Default to internal server error
	return fiber.StatusInternalServerError
}

// NewRecoveryMiddleware provides panic recovery
func NewRecoveryMiddleware() fiber.Handler {
	recoveryLogger, err := logger.New("recovery")
	if err != nil {
		panic("failed to create recovery logger: " + err.Error())
	}

	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				requestID := GetRequestID(c)

				// Log the panic
				recoveryLogger.Error("panic recovered",
					zap.String("panic", fmt.Sprintf("%v", r)),
					zap.String("request_id", requestID),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.String("ip", GetRealIP(c)),
				)

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
				err := c.Status(fiber.StatusInternalServerError).JSON(response)
				if err != nil {
					return
				}
			}
		}()

		return c.Next()
	}
}
