package middleware

import (
	"net/http"

	"github.com/adampresley/logging"
	"github.com/gorilla/context"
	"github.com/mailslurper/libmailslurper/storage"
)

/*
AppContext holds context data for the application. This can hold information
such as a database connection, session data, user info, and more. Your middlewares
should attach functions to this structure to pass critical data to request
handlers.
*/
type AppContext struct {
	Log      *logging.Logger
	Database storage.IStorage
}

/*
StartAppContext is a middleware that should be early in the chain. This
sets up the initial context and attaches important data to the Gorilla
Context which comes across in the request.
*/
func (ctx *AppContext) StartAppContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		context.Set(request, "log", ctx.Log)
		context.Set(request, "database", ctx.Database)

		h.ServeHTTP(writer, request)

		context.Clear(request)
	})
}
