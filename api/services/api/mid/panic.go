package mid

import (
	"context"
	"net/http"

	"github.com/Farber98/garage-sale/app/api/mid"
	"github.com/Farber98/garage-sale/foundation/web"
)

// Panic executes the panic middleware functionality.
func Panic() web.MidHandler {
	m := func(handler web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Panics(ctx, hdl)
		}
		return h
	}
	return m
}
