package logging

import "log/slog"

// Logger abtracts any logger used by internal modules.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}

// NewLoggerFromSlog creates a new [Logger] using [slog.Logger].
func NewLoggerFromSlog(logger *slog.Logger) Logger {
	return &slogLogger{
		log: logger,
	}
}

// slogLogger wraps [slog.Logger] to implement [Logger] interface.
type slogLogger struct {
	log *slog.Logger
}

// Debug implements Logger.
func (s *slogLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, args...)
}

// Error implements Logger.
func (s *slogLogger) Error(msg string, args ...any) {
	s.log.Error(msg, args...)
}

// Info implements Logger.
func (s *slogLogger) Info(msg string, args ...any) {
	s.log.Info(msg, args...)
}

// Warn implements Logger.
func (s *slogLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, args...)
}

// With implements Logger.
func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		log: s.log.With(args...),
	}
}
