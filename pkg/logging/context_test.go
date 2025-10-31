package logging

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithContext(t *testing.T) {
	t.Run("should store the logger as context value", func(t *testing.T) {
		// given
		ctx := t.Context()
		logger := &slogLogger{}

		// when
		newCtx := WithContext(ctx, logger)

		// then
		assert.Same(t, logger, newCtx.Value(ctxKey("logger")))
	})
}

func TestFromContext(t *testing.T) {
	t.Run("should return the logger value from the context", func(t *testing.T) {
		// given
		expectedLogger := &slogLogger{}
		ctx := context.WithValue(t.Context(), ctxKey("logger"), expectedLogger)

		// when
		logger := FromContext(ctx)

		// then
		assert.Same(t, expectedLogger, logger)
	})

	t.Run("should return new slog logger if no logger has been added to the context", func(t *testing.T) {
		// given
		ctx := t.Context()

		// when
		logger := FromContext(ctx)

		// then
		assert.Equal(t, &slogLogger{
			log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		}, logger)
	})
}
