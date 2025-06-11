package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		want     string
	}{
		{
			name: "error with cause",
			appError: &AppError{
				Type:    "TestError",
				Message: "test message",
				Cause:   errors.New("underlying error"),
			},
			want: "TestError: test message: underlying error",
		},
		{
			name: "error without cause",
			appError: &AppError{
				Type:    "TestError",
				Message: "test message",
				Cause:   nil,
			},
			want: "TestError: test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appError.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	appError := &AppError{
		Type:    "TestError",
		Message: "test message",
		Cause:   cause,
	}

	if got := appError.Unwrap(); got != cause {
		t.Errorf("AppError.Unwrap() = %v, want %v", got, cause)
	}

	// Test with no cause
	appErrorNoCause := &AppError{
		Type:    "TestError",
		Message: "test message",
		Cause:   nil,
	}

	if got := appErrorNoCause.Unwrap(); got != nil {
		t.Errorf("AppError.Unwrap() = %v, want nil", got)
	}
}

func TestNewAppError(t *testing.T) {
	cause := errors.New("underlying error")
	appError := NewAppError("TestError", "test message", cause)

	if appError.Type != "TestError" {
		t.Errorf("NewAppError() Type = %v, want %v", appError.Type, "TestError")
	}
	if appError.Message != "test message" {
		t.Errorf("NewAppError() Message = %v, want %v", appError.Message, "test message")
	}
	if appError.Cause != cause {
		t.Errorf("NewAppError() Cause = %v, want %v", appError.Cause, cause)
	}
}

func TestNewConfigError(t *testing.T) {
	cause := errors.New("config error")
	configError := NewConfigError("invalid config", cause)

	if configError.Type != "ConfigError" {
		t.Errorf("NewConfigError() Type = %v, want %v", configError.Type, "ConfigError")
	}
	if configError.Message != "invalid config" {
		t.Errorf("NewConfigError() Message = %v, want %v", configError.Message, "invalid config")
	}
	if configError.Cause != cause {
		t.Errorf("NewConfigError() Cause = %v, want %v", configError.Cause, cause)
	}
}

func TestNewDatabaseError(t *testing.T) {
	cause := errors.New("db connection failed")
	dbError := NewDatabaseError("database error", cause)

	if dbError.Type != "DatabaseError" {
		t.Errorf("NewDatabaseError() Type = %v, want %v", dbError.Type, "DatabaseError")
	}
	if dbError.Message != "database error" {
		t.Errorf("NewDatabaseError() Message = %v, want %v", dbError.Message, "database error")
	}
	if dbError.Cause != cause {
		t.Errorf("NewDatabaseError() Cause = %v, want %v", dbError.Cause, cause)
	}
}

func TestNewNotFoundError(t *testing.T) {
	notFoundError := NewNotFoundError("user")

	if notFoundError.Type != "NotFoundError" {
		t.Errorf("NewNotFoundError() Type = %v, want %v", notFoundError.Type, "NotFoundError")
	}
	if notFoundError.Message != "user not found" {
		t.Errorf("NewNotFoundError() Message = %v, want %v", notFoundError.Message, "user not found")
	}
	if notFoundError.Code != 404 {
		t.Errorf("NewNotFoundError() Code = %v, want %v", notFoundError.Code, 404)
	}
}

func TestNewInternalError(t *testing.T) {
	cause := errors.New("internal error")
	internalError := NewInternalError("something went wrong", cause)

	if internalError.Type != "InternalError" {
		t.Errorf("NewInternalError() Type = %v, want %v", internalError.Type, "InternalError")
	}
	if internalError.Message != "something went wrong" {
		t.Errorf("NewInternalError() Message = %v, want %v", internalError.Message, "something went wrong")
	}
	if internalError.Code != 500 {
		t.Errorf("NewInternalError() Code = %v, want %v", internalError.Code, 500)
	}
	if internalError.Cause != cause {
		t.Errorf("NewInternalError() Cause = %v, want %v", internalError.Cause, cause)
	}
}

func TestIsAppError(t *testing.T) {
	appError := NewAppError("TestError", "test message", nil)
	regularError := errors.New("regular error")

	if !IsAppError(appError) {
		t.Errorf("IsAppError() should return true for AppError")
	}

	if IsAppError(regularError) {
		t.Errorf("IsAppError() should return false for regular error")
	}

	if IsAppError(nil) {
		t.Errorf("IsAppError() should return false for nil")
	}
}

func TestGetAppError(t *testing.T) {
	appError := NewAppError("TestError", "test message", nil)
	regularError := errors.New("regular error")

	// Test with AppError
	if got := GetAppError(appError); got != appError {
		t.Errorf("GetAppError() = %v, want %v", got, appError)
	}

	// Test with regular error
	if got := GetAppError(regularError); got != nil {
		t.Errorf("GetAppError() = %v, want nil", got)
	}

	// Test with wrapped AppError
	wrappedError := fmt.Errorf("wrapped: %w", appError)
	if got := GetAppError(wrappedError); got != appError {
		t.Errorf("GetAppError() = %v, want %v", got, appError)
	}
}

func TestWrap(t *testing.T) {
	originalError := errors.New("original error")
	wrappedError := Wrap(originalError, "additional context")

	expectedMessage := "additional context: original error"
	if wrappedError.Error() != expectedMessage {
		t.Errorf("Wrap() = %v, want %v", wrappedError.Error(), expectedMessage)
	}

	// Test wrapping nil error
	if got := Wrap(nil, "context"); got != nil {
		t.Errorf("Wrap(nil, context) = %v, want nil", got)
	}
}

func TestWrapf(t *testing.T) {
	originalError := errors.New("original error")
	wrappedError := Wrapf(originalError, "context with %s", "formatting")

	expectedMessage := "context with formatting: original error"
	if wrappedError.Error() != expectedMessage {
		t.Errorf("Wrapf() = %v, want %v", wrappedError.Error(), expectedMessage)
	}

	// Test wrapping nil error
	if got := Wrapf(nil, "context with %s", "formatting"); got != nil {
		t.Errorf("Wrapf(nil, context) = %v, want nil", got)
	}
}
