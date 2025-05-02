package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

func (app *App) HandleFunc(pattern string, handler Handler) {

	// Callback func to inject custom logic
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
		}
	}

	app.ServeMux.HandleFunc(pattern, h)
}
