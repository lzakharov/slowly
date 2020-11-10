package slowly

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lzakharov/slowly/pkg/web"
)

// App is a Slowly application.
type App struct {
	srv *http.Server
}

// NewApp creates a new application.
func NewApp(addr string, timeout time.Duration) *App {
	withTimeout := web.WithTimeout(timeout, func(w http.ResponseWriter) {
		web.WriteJSON(w, http.StatusBadRequest, TimeoutTooLongResponse)
	})

	mux := http.NewServeMux()
	mux.Handle("/api/slow", withTimeout(http.HandlerFunc(Slow)))

	return &App{srv: &http.Server{Addr: addr, Handler: mux}}
}

// Run runs the application.
func (a App) Run() error {
	err := a.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the application.
func (a App) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
