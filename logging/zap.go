package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type zapLogger struct {
	*zap.SugaredLogger
}

// NewZap creates a new IshiLogger backed by a zap SugaredLogger.
func NewZap(env string) (IshiLogger, error) {
	var l *zap.Logger
	switch env {
	case "local", "dev":
		option := zap.AddStacktrace(zap.ErrorLevel)
		l, _ = zap.NewDevelopment(option)
	case "qa", "uat", "prod":
		l, _ = zap.NewProduction()
	default:
		return nil, fmt.Errorf("invalid environment specified: %q", env)
	}
	return zapLogger{l.Sugar()}, nil
}

// Debug provides structured logging at the debug level.
func (zl zapLogger) Debug(msg string, args ...interface{}) {
	zl.Debugw(msg, args)
}

// Info provides structured logging at the info level.
func (zl zapLogger) Info(msg string, args ...interface{}) {
	zl.Infow(msg, args)
}

// Warn provides structured logging at the warn level.
func (zl zapLogger) Warn(msg string, args ...interface{}) {
	zl.Warnw(msg, args)
}

// Error provides structured logging at the error level.
func (zl zapLogger) Error(msg string, args ...interface{}) {
	zl.Errorw(msg, args)
}

// Fatal provides structured logging at the fatal level.
func (zl zapLogger) Fatal(msg string, args ...interface{}) {
	zl.Fatalw(msg, args)
}

// Panic provides structured logging at the panic level, panicking afterward.
func (zl zapLogger) Panic(msg string, args ...interface{}) {
	zl.Panicw(msg, args)
}

// Sync flushes any buffered log entries.
func (zl zapLogger) Sync() error {
	return zl.Sync()
}

// With adds a variadic number of fields to the logging context.
func (zl zapLogger) WithFields(args ...interface{}) IshiLogger {
	return zapLogger{zl.With(args)}
}
