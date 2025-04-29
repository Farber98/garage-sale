package mux

import (
	"net/http"

	"github.com/Farber98/garage-sale/api/services/sales/route/sys/checkapi"
)

// WebAPI constructs a http.Handler with all app routes bound
func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	checkapi.Routes(mux)

	return mux
}
