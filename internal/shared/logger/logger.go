package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
}

type contextKey string

const requestIDKey contextKey = "request_id"

func New(level string) *Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &Logger{Logger: logger}
}

func NewWithWriter(w io.Writer, level string) *Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(w, opts)
	return &Logger{Logger: slog.New(handler)}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With(slog.String("request_id", requestID))}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return l.WithRequestID(requestID)
	}
	return l
}

func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{Logger: l.Logger.With(slog.Any(key, value))}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With(slog.String("error", err.Error()))}
}

func (l *Logger) WithDuration(d time.Duration) *Logger {
	return &Logger{Logger: l.Logger.With(slog.Duration("duration", d))}
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelDebug, msg, args...)
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelInfo, msg, args...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelWarn, msg, args...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelError, msg, args...)
}

func RequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}