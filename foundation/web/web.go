package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

func (app *App) SignalShutdown() {
	app.shutdown <- syscall.SIGTERM
}

func (app *App) HandleFunc(pattern string, handler HandlerFunc, mw ...MidHandler) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(app.mw, handler)

	// Callback func to inject custom logic
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}

		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				app.SignalShutdown()
				return
			}
		}
	}

	app.ServeMux.HandleFunc(pattern, h)
}

func (app *App) HandleFuncNoMiddleware(pattern string, handler HandlerFunc, mw ...MidHandler) {
	// Callback func to inject custom logic
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}

		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				app.SignalShutdown()
				return
			}
		}
	}

	app.ServeMux.HandleFunc(pattern, h)
}

func validateError(err error) bool {
	switch {
	case errors.Is(err, syscall.EPIPE):
		return false
	case errors.Is(err, syscall.ECONNRESET):
		return false
	}

	return true
}
