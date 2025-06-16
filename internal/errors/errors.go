package errors

import (
	"errors"
	"fmt"
)

// AppError represents an application-specific error with context
type AppError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
	Code    int    `json:"code,omitempty"`
}

func New(text string) error {
	return errors.New(text)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error
func NewAppError(errorType, message string, cause error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigError creates a configuration error
func NewConfigError(message string, cause error) *AppError {
	return &AppError{
		Type:    "ConfigError",
		Message: message,
		Cause:   cause,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string, cause error) *AppError {
	return &AppError{
		Type:    "DatabaseError",
		Message: message,
		Cause:   cause,
	}
}

// NewElasticsearchError creates an Elasticsearch error
func NewElasticsearchError(message string, cause error) *AppError {
	return &AppError{
		Type:    "ElasticsearchError",
		Message: message,
		Cause:   cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    "ValidationError",
		Message: message,
		Cause:   cause,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:    "NotFoundError",
		Message: fmt.Sprintf("%s not found", resource),
		Code:    404,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    "InternalError",
		Message: message,
		Cause:   cause,
		Code:    500,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from error chain
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
