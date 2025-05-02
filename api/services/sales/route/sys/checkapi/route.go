package checkapi

import (
	"github.com/Farber98/garage-sale/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /readiness", readiness)
	app.HandleFunc("GET /liveness", liveness)
}
