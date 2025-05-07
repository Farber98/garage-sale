package checkapi

import (
	"github.com/Farber98/garage-sale/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFunc("GET /testerror", testerror)
	app.HandleFunc("GET /testpanic", testpanic)
}
