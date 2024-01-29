package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
}

func New() *Router {
	return &Router{}
}

func (router *Router) Routes(r chi.Router) {
	r.Group(router.test)
	r.Group(router.users)
}

// groups

func (router *Router) users(r chi.Router) {

}

func (router *Router) test(r chi.Router) {
	r.Use(middleware.Logger)
	const prefix = "/api/v1/ping"
	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}
