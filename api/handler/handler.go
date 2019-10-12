package handler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Route is a contract to bind our http routers
type Route interface {
	RegisterHandler(r *httprouter.Router) error
}

// NewHandler return standard handlers for our service
func NewHandler(routes ...Route) http.Handler {
	router := httprouter.New()

	router.HandlerFunc("GET", "/healthz", Healthz)
	for _, r := range routes {
		r.RegisterHandler(router)
	}

	return router
}

// Healthz is our service healthiness handler
func Healthz(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "ok")
}

// Decorate applies middlewares to our base handler
func Decorate(handler StandardHandler, middlewares ...Middleware) httprouter.Handle {
	return HTTP(AppendMiddlewares(handler, middlewares...))
}
