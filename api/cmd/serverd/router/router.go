package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	websocketHandler "github.com/tanmaij/zylo/internal/handler/ws"
	"github.com/tanmaij/zylo/pkg/ws"
)

type Router struct {
	WSHandler  websocketHandler.Handler
	WSListener *ws.WebsocketListener
}

func New(wsListener *ws.WebsocketListener, wsHandler websocketHandler.Handler) *Router {
	return &Router{WSHandler: wsHandler, WSListener: wsListener}
}

func (router *Router) Routes(r chi.Router) {
	r.Group(router.test)
	r.Group(router.websocket)
	r.Group(router.users)
}

// groups

func (router *Router) users(r chi.Router) {

}

func (router *Router) websocket(r chi.Router) {
	router.WSListener.OnConnection(router.WSHandler.OnConnection)
	router.WSListener.OnDisconnection(router.WSHandler.OnDisconnection)
	router.WSListener.On("ping", router.WSHandler.Ping)

	r.Get("/ws", router.WSListener.HandleConnections)
}

func (router *Router) test(r chi.Router) {
	r.Use(middleware.Logger)
	const prefix = "/api/v1/ping"
	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}
