package mux

import (
	"os"

	"github.com/Farber98/garage-sale/api/services/sales/route/sys/checkapi"
	"github.com/Farber98/garage-sale/foundation/web"
)

// WebAPI constructs a http.Handler with all app routes bound
func WebAPI(shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown)

	checkapi.Routes(app)

	return app
}
