package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide a consistent interface
type Logger struct {
	zap    *zap.Logger
	sugar  *zap.SugaredLogger
	prefix string
}

// Fields represents structured log fields
type Fields map[string]any

var (
	defaultLogger *Logger
	globalConfig  zap.Config
)

// init initializes the default logger
func init() {
	var err error
	defaultLogger, err = New("")
	if err != nil {
		panic("failed to initialize default logger: " + err.Error())
	}
}

// New creates a new logger with the specified prefix
func New(prefix string) (*Logger, error) {
	config := getLoggerConfig()

	zapLogger, err := config.Build(
		zap.AddCallerSkip(1), // 跳过封装层
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		zap:    zapLogger,
		sugar:  zapLogger.Sugar(),
		prefix: prefix,
	}, nil
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// getLoggerConfig returns the zap configuration
func getLoggerConfig() zap.Config {
	config := zap.NewProductionConfig()

	// Customize the configuration
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Development = false
	config.DisableCaller = false
	config.DisableStacktrace = false
	config.Sampling = nil // Disable sampling for now

	// Customize the encoder
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use console encoder for better readability
	config.Encoding = "console"

	return config
}

// SetLevel sets the minimum log level for the logger
func (l *Logger) SetLevel(level zapcore.Level) {
	if l.zap != nil {
		l.zap = l.zap.WithOptions(zap.IncreaseLevel(level))
		l.sugar = l.zap.Sugar()
	}
}

// SetGlobalLevel sets the global log level
func SetGlobalLevel(level zapcore.Level) {
	if defaultLogger != nil {
		defaultLogger.SetLevel(level)
	}
}

// fieldsToZapFields converts Fields to zap fields
func fieldsToZapFields(fields Fields) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// Close closes the logger and flushes any buffered log entries
func (l *Logger) Close() error {
	if l.zap != nil {
		return l.zap.Sync()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Fields) {
	if l.zap == nil {
		return
	}
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.zap.Debug(msg, fieldsToZapFields(f)...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Fields) {
	if l.zap == nil {
		return
	}
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.zap.Info(msg, fieldsToZapFields(f)...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Fields) {
	if l.zap == nil {
		return
	}
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.zap.Warn(msg, fieldsToZapFields(f)...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Fields) {
	if l.zap == nil {
		return
	}
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.zap.Error(msg, fieldsToZapFields(f)...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...Fields) {
	if l.zap == nil {
		os.Exit(1)
		return
	}
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.zap.Fatal(msg, fieldsToZapFields(f)...)
}

// WithContext adds context information to fields
func WithContext(ctx context.Context) Fields {
	fields := make(Fields)

	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}

	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}

	return fields
}

// WithError adds error information to fields
func WithError(err error) Fields {
	return Fields{
		"error": err.Error(),
	}
}

// Debug Global logger functions
func Debug(msg string, fields ...Fields) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...Fields) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...Fields) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...Fields) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...Fields) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	} else {
		os.Exit(1)
	}
}

// Sync flushes any buffered log entries for the default logger
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}
