package app

import (
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestSignalContext(t *testing.T) {
	t.Run("should cancel context on signals", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			signal syscall.Signal
		}{
			{
				name:   "SIGINT",
				signal: syscall.SIGINT,
			},
			{
				name:   "SIGTERM",
				signal: syscall.SIGTERM,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// given
				signal.Reset()
				ctx := SignalContext()

				// when
				syscall.Kill(syscall.Getpid(), tt.signal)

				// then
				select {
				case <-ctx.Done():
				case <-time.After(time.Second):
					t.Error("context not canceled after one second")
				}
			})
		}
	})
}
