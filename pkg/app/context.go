package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalContext returns a new context and handles SIGINT and SIGTERM signals to
// cancel it. This context can be used to implement graceful shutdowns.
func SignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()

		<-sigChan
		os.Exit(1)
	}()

	return ctx
}
