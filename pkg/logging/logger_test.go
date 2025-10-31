package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerFromSlog(t *testing.T) {
	t.Run("should return a logger based on slog logger", func(t *testing.T) {
		// given
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		// when
		newLogger := NewLoggerFromSlog(logger)

		// then
		assert.Equal(t, &slogLogger{
			log: logger,
		}, newLogger)
	})
}

func TestSlogLogger_Level(t *testing.T) {
	for _, tt := range []struct {
		name          string
		logFn         func(logger *slogLogger) func(string, ...any)
		expectedLevel string
	}{
		{
			name:          "debug",
			logFn:         func(logger *slogLogger) func(string, ...any) { return logger.Debug },
			expectedLevel: "DEBUG",
		},
		{
			name:          "info",
			logFn:         func(logger *slogLogger) func(string, ...any) { return logger.Info },
			expectedLevel: "INFO",
		},
		{
			name:          "warn",
			logFn:         func(logger *slogLogger) func(string, ...any) { return logger.Warn },
			expectedLevel: "WARN",
		},
		{
			name:          "error",
			logFn:         func(logger *slogLogger) func(string, ...any) { return logger.Error },
			expectedLevel: "ERROR",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// given
			var buf bytes.Buffer
			logger := &slogLogger{log: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))}

			// when
			tt.logFn(logger)("test message", "key", "value")

			// then
			var m map[string]any
			assert.NoError(t, json.Unmarshal(buf.Bytes(), &m))
			assert.NotNil(t, m["time"])
			assert.Equal(t, "test message", m["msg"])
			assert.Equal(t, tt.expectedLevel, m["level"])
			assert.Equal(t, "value", m["key"])
		})
	}
}

func TestSlogLogger_With(t *testing.T) {
	// given
	var buf bytes.Buffer
	logger := &slogLogger{log: slog.New(slog.NewJSONHandler(&buf, nil))}

	// when
	newLogger := logger.With("key", "value")
	newLogger.Error("some error")

	// then
	var m map[string]any
	assert.NotSame(t, logger, newLogger)
	assert.NoError(t, json.Unmarshal(buf.Bytes(), &m))
	assert.Equal(t, "some error", m["msg"])
	assert.Equal(t, "value", m["key"])
}
