// Package logging provides structured logging utilities using zap.
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with project-specific helpers.
type Logger struct{ zap *zap.Logger }

// New creates a production JSON logger at the given level (debug/info/warn/error).
func New(level string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	z, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{zap: z}, nil
}

// Info logs at info level.
func (l *Logger) Info(msg string, fields ...zap.Field) { l.zap.Info(msg, fields...) }

// Error logs at error level.
func (l *Logger) Error(msg string, fields ...zap.Field) { l.zap.Error(msg, fields...) }

// Warn logs at warn level.
func (l *Logger) Warn(msg string, fields ...zap.Field) { l.zap.Warn(msg, fields...) }

// Debug logs at debug level.
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.zap.Debug(msg, fields...) }

// With returns a new Logger with the given fields.
func (l *Logger) With(fields ...zap.Field) *Logger { return &Logger{zap: l.zap.With(fields...)} }

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error { return l.zap.Sync() }
