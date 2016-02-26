package libmailslurper

import (
	"github.com/adampresley/logging"
	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/libmailslurper/middleware"
	"github.com/mailslurper/libmailslurper/server"
)

/*
StartServiceTier starts the service tier HTTP application. This will
setup logging, an app context, and the HTTP listener
*/
func StartServiceTier(configuration *configuration.ServiceTierConfiguration) error {
	log := logging.NewLoggerWithMinimumLevel("MailSlurper", logging.StringToLogType("info"))
	appContext := &middleware.AppContext{
		Log:      log,
		Database: configuration.Database,
	}

	httpListener := server.NewHTTPListenerService(configuration.Address, configuration.Port, appContext)

	setupMiddleware(httpListener, appContext)
	setupRoutes(httpListener, appContext)

	return httpListener.StartHTTPListener()
}
