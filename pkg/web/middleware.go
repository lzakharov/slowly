package web

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Middleware is a HTTP middleware.
type Middleware func(http.Handler) http.Handler

// WithTimeout creates a middleware that injects context with a given timeout and runs callback on the deadline.
func WithTimeout(timeout time.Duration, callback func(http.ResponseWriter)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					callback(w)
				}
			}()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
