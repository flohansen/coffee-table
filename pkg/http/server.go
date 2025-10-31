package http

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Mux is a [http.Handler] with additional functions to implement routing
// requests.
type Mux interface {
	http.Handler
	Handle(pattern string, handler http.Handler)
}

// Middleware wraps a [http.Handler] to implement pre- or post-processing of a
// request/response.
type Middleware func(http.Handler) http.Handler

// Server is an HTTP server handling incoming requests.
type Server struct {
	addr        string
	mux         Mux
	middlewares []Middleware
}

// NewServer returns a new Server, which is configured using the given options.
// The default listen address is ":8080".
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		addr: ":8080",
		mux:  http.NewServeMux(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Serve starts the HTTP server on a real network (TCP) address. Before the
// server is being started, the middlewares are applied to the initial handler
// (request multiplexer) with respect to their order. If the context is being
// canceled, the server will be shutdown gracefully.
func (s *Server) Serve(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.applyMiddlewares(),
	}

	srvErr := make(chan error, 1)
	go func() {
		defer close(srvErr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
		return nil
	case err := <-srvErr:
		return err
	}
}

// applyMiddlewares applies the registered middlewares with respect to their
// order.
func (s *Server) applyMiddlewares() http.Handler {
	handler := http.Handler(s.mux)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}
	return handler
}

// ServerOption configures a Server.
type ServerOption func(*Server)

// WithListenAddr sets the TCP listen address of a Server.
func WithListenAddr(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithHandler registeres the HTTP handler for a pattern to route HTTP
// requests.
func WithHandler(pattern string, handler http.Handler) ServerOption {
	return func(s *Server) {
		s.mux.Handle(pattern, handler)
	}
}

// WithMiddleware adds the middlewares to a Server. They are not applied
// immediately, but only when starting the server with [Serve].
func WithMiddleware(middlewares ...Middleware) ServerOption {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, middlewares...)
	}
}

// WithMux sets the request multiplexer of a Server to route incoming HTTP
// requests.
func WithMux(mux Mux) ServerOption {
	return func(s *Server) {
		s.mux = mux
	}
}
