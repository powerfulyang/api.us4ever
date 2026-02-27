package logger

import (
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

func IsLocalDev(appEnv string) bool {
	return appEnv == "local"
}

// New creates a new logger with the specified prefix
func New(prefix string) (*Logger, error) {
	config := getLoggerConfig()

	appEnv := os.Getenv("APP_ENV")
	if IsLocalDev(appEnv) {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

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
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
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

// Close closes the logger and flushes any buffered log entries
func (l *Logger) Close() error {
	if l.zap != nil {
		return l.zap.Sync()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if l.zap == nil {
		return
	}
	l.zap.Debug(msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	if l.zap == nil {
		return
	}
	l.zap.Info(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if l.zap == nil {
		return
	}
	l.zap.Warn(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	if l.zap == nil {
		return
	}
	l.zap.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if l.zap == nil {
		os.Exit(1)
		return
	}
	l.zap.Fatal(msg, fields...)
}

// Debugw logs a debug message with key-value pairs (sugar style)
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	if l.sugar == nil {
		return
	}
	l.sugar.Debugw(msg, keysAndValues...)
}

// Infow logs an info message with key-value pairs (sugar style)
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	if l.sugar == nil {
		return
	}
	l.sugar.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with key-value pairs (sugar style)
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	if l.sugar == nil {
		return
	}
	l.sugar.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message with key-value pairs (sugar style)
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	if l.sugar == nil {
		return
	}
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatalw logs a fatal message with key-value pairs (sugar style) and exits
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	if l.sugar == nil {
		os.Exit(1)
		return
	}
	l.sugar.Fatalw(msg, keysAndValues...)
}
