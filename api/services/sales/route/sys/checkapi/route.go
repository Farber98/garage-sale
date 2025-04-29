package checkapi

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /readiness", readiness)
	mux.HandleFunc("GET /liveness", liveness)
}
