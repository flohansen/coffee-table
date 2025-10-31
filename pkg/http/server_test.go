package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_NewServer(t *testing.T) {
	t.Run("default listen address", func(t *testing.T) {
		// given
		// when
		srv := NewServer(WithListenAddr(":8080"))

		// then
		assert.Equal(t, ":8080", srv.addr)
	})

	t.Run("options", func(t *testing.T) {
		t.Run("with listen address", func(t *testing.T) {
			// given
			// when
			srv := NewServer(WithListenAddr(":8888"))

			// then
			assert.Equal(t, ":8888", srv.addr)
		})

		t.Run("with handler", func(t *testing.T) {
			// given
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			// when
			srv := NewServer(
				WithHandler("/", handler),
			)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			srv.mux.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, 200, res.StatusCode)
		})

		t.Run("with middlewares", func(t *testing.T) {
			// given
			// when
			srv := NewServer(
				WithMiddleware(
					writeMiddleware(""),
					writeMiddleware(""),
					writeMiddleware(""),
				),
			)

			// then
			assert.Len(t, srv.middlewares, 3)
		})

		t.Run("with mux", func(t *testing.T) {
			// given
			mux := http.NewServeMux()

			// when
			srv := NewServer(WithMux(mux))

			// then
			assert.Same(t, mux, srv.mux)
		})
	})
}

func TestServer_Serve(t *testing.T) {
	t.Run("should shutdown gracefully", func(t *testing.T) {
		// given
		srv := NewServer()
		ctx, cancel := context.WithCancel(t.Context())

		// when
		cancel()
		err := srv.Serve(ctx)

		// then
		assert.NoError(t, err)
	})

	t.Run("should apply middlewares in correct order", func(t *testing.T) {
		// given
		srv := NewServer(
			WithHandler("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test"))
			})),
			WithMiddleware(
				writeMiddleware("this"),
				writeMiddleware("is"),
				writeMiddleware("a"),
			),
		)

		// when
		handler := srv.applyMiddlewares()
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(w, r)

		// then
		res := w.Result()
		b, _ := io.ReadAll(res.Body)
		assert.Equal(t, "thisisatest", string(b))
	})
}

func writeMiddleware(str string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(str))
			next.ServeHTTP(w, r)
		})
	}
}
