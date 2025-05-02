package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

func (app *App) HandleFunc(pattern string, handler HandlerFunc, mw ...MidHandler) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(app.mw, handler)

	// Callback func to inject custom logic
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
		}
	}

	app.ServeMux.HandleFunc(pattern, h)
}
