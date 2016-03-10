package server

import (
	"fmt"
	"net/http"

	"github.com/adampresley/logging"
	"github.com/mailslurper/libmailslurper/middleware"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

/*
HTTPListenerService is a structure which provides an HTTP listener to service
requests. This structure offers methods to add routes and middlewares. Typical
usage would first call NewHTTPListenerService(), add routes, then call
StartHTTPListener.
*/
type HTTPListenerService struct {
	Address string
	Port    int
	Context *middleware.AppContext

	Router                 *mux.Router
	BaseMiddlewareHandlers alice.Chain
}

/*
NewHTTPListenerService creates a new instance of the HTTPListenerService
*/
func NewHTTPListenerService(
	address string,
	port int,
	appContext *middleware.AppContext,
) *HTTPListenerService {
	return &HTTPListenerService{
		Address: address,
		Port:    port,
		Context: appContext,

		Router: mux.NewRouter(),
	}
}

/*
AddMiddleware adds a new middleware handler to the request chain.
*/
func (service *HTTPListenerService) AddMiddleware(middlewareHandler alice.Constructor) *HTTPListenerService {
	service.BaseMiddlewareHandlers = service.BaseMiddlewareHandlers.Append(middlewareHandler)
	return service
}

/*
AddRoute adds a HTTP handler route to the HTTP listener.
*/
func (service *HTTPListenerService) AddRoute(
	path string,
	handlerFunc http.HandlerFunc, methods ...string,
) *HTTPListenerService {
	service.Router.Handle(path, service.BaseMiddlewareHandlers.ThenFunc(handlerFunc)).Methods(methods...)
	return service
}

/*
AddRouteWithMiddleware adds a HTTP handler route that goes through an additional
middleware handler, to the HTTP listener.
*/
func (service *HTTPListenerService) AddRouteWithMiddleware(
	path string,
	handlerFunc http.HandlerFunc,
	middlewareHandler alice.Constructor,
	methods ...string,
) *HTTPListenerService {
	service.Router.Handle(
		path,
		service.BaseMiddlewareHandlers.Append(middlewareHandler).ThenFunc(handlerFunc),
	).Methods(methods...)

	return service
}

/*
AddStaticRoute adds a HTTP handler route for static assets.
*/
func (service *HTTPListenerService) AddStaticRoute(pathPrefix string, directory string) *HTTPListenerService {
	fileServer := http.FileServer(http.Dir("./www/assets"))
	service.Router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, fileServer))

	return service
}

/*
StartHTTPListener starts the HTTP listener and servicing requests.
*/
func (service *HTTPListenerService) StartHTTPListener() error {
	listener := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", service.Address, service.Port),
		Handler: alice.New().Then(service.Router),
	}

	return startListener(listener, service.Context.CertFile, service.Context.KeyFile, service.Context.Log, service.Address, service.Port)
}

func startListener(listener *http.Server, certFile, keyFile string, logger *logging.Logger, address string, port int) error {
	if certFile != "" && keyFile != "" {
		logger.Infof("Service Tier HTTPS listener started on %s:%d", address, port)
		return listener.ListenAndServeTLS(certFile, keyFile)
	}

	logger.Infof("Service Tier HTTP listener started on %s:%d", address, port)
	return listener.ListenAndServe()
}
