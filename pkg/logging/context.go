package logging

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	ctxKeyLogger ctxKey = "logger"
)

// WithContext adds a logger as context value to a parent context.
func WithContext(parent context.Context, logger Logger) context.Context {
	return context.WithValue(parent, ctxKeyLogger, logger)
}

// FromContext returns a logger stored as context value from a context. If
// there is no logger, a new [slog.Logger] will be returned.
func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(ctxKeyLogger).(Logger)
	if !ok {
		return NewLoggerFromSlog(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	return logger
}
