package libmailslurper

import (
	"github.com/mailslurper/libmailslurper/controllers"
	"github.com/mailslurper/libmailslurper/middleware"
	"github.com/mailslurper/libmailslurper/server"
)

/*
Add routes here using AddRoute and AddRouteWithMiddleware.
*/
func setupRoutes(httpListener *server.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddRoute("/version", controllers.GetVersion, "GET", "OPTIONS").
		AddRoute("/mail/{mailID}", controllers.GetMail, "GET", "OPTIONS").
		AddRoute("/mail/{mailID}/message", controllers.GetMailMessage, "GET", "OPTIONS").
		AddRoute("/mail/{mailID}/attachment/{attachmentID}", controllers.DownloadAttachment, "GET", "OPTIONS").
		AddRoute("/mail", controllers.GetMailCollection, "GET", "OPTIONS").
		AddRoute("/mail", controllers.DeleteMail, "DELETE", "OPTIONS").
		AddRoute("/mailcount", controllers.GetMailCount, "GET", "OPTIONS").
		AddRoute("/pruneoptions", controllers.GetPruneOptions, "GET", "OPTIONS")
}
