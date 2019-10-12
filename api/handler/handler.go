package handler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	dms = DefaultMiddlewares()

	AppAuth  []Middleware
	UserAuth []Middleware
)

// Route is a contract to bind our http routers
type Route interface {
	RegisterHandler(r *httprouter.Router) error
}

func initAuthMiddlewares(privateKey string) {
	UserAuth = append([]Middleware{WithAuthentication(privateKey)}, dms...)
	AppAuth = dms
}

// NewHandler return standard handlers for our service
func NewHandler(jwtPrivateKey string, routes ...Route) http.Handler {
	initAuthMiddlewares(jwtPrivateKey)

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
