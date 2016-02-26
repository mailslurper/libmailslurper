package middleware

import (
	"net/http"
	"net/url"
	"time"
)

/*
Logger is a middleware which logs requests to the console. It also includes the
time it takes for the request to complete.
*/
func (ctx *AppContext) Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(writer, request)

		requestURL, _ := url.QueryUnescape(request.URL.String())
		ctx.Log.Infof("%s - %s (%s)", request.Method, requestURL, time.Since(startTime))
	})
}
