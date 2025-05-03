package mux

import (
	"os"

	"github.com/Farber98/garage-sale/api/services/sales/route/sys/checkapi"
	"github.com/Farber98/garage-sale/app/mid"
	"github.com/Farber98/garage-sale/foundation/logger"
	"github.com/Farber98/garage-sale/foundation/web"
)

// WebAPI constructs a http.Handler with all app routes bound
func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log))

	checkapi.Routes(app)

	return app
}
