package libmailslurper

import (
	"github.com/mailslurper/libmailslurper/middleware"
	"github.com/mailslurper/libmailslurper/server"
)

func setupMiddleware(httpListener *server.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddMiddleware(appContext.Logger).
		AddMiddleware(appContext.StartAppContext).
		AddMiddleware(appContext.AccessControl).
		AddMiddleware(appContext.OptionsHandler)
}
