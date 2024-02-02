package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	websocketHandler "github.com/tanmaij/zylo/internal/handler/ws"
	"github.com/tanmaij/zylo/pkg/ws"
)

// Router struct represents the main router for handling HTTP and WebSocket routes.
type Router struct {
	WSHandler  websocketHandler.Handler
	WSListener *ws.WebsocketListener
}

// New creates a new instance of the Router.
func New(wsListener *ws.WebsocketListener, wsHandler websocketHandler.Handler) *Router {
	return &Router{WSHandler: wsHandler, WSListener: wsListener}
}

// RegisterRoutes registers different route groups.
func (router *Router) RegisterRoutes(r chi.Router) {
	r.Group(router.test)
	r.Group(router.websocket)
	r.Group(router.users)
}

// users is a route group for handling user-related routes.
func (router *Router) users(r chi.Router) {
	// Add user-related route handlers here
}

// websocket is a route group for handling WebSocket connections.
func (router *Router) websocket(r chi.Router) {
	// Set up WebSocket connection handling using WSListener and WSHandler
	router.WSListener.OnConnection(router.WSHandler.OnConnection)
	router.WSListener.OnDisconnection(router.WSHandler.OnDisconnection)
	router.WSListener.On("ping", router.WSHandler.Ping)

	// Route for WebSocket connection endpoint
	r.Get("/ws", router.WSListener.HandleConnections)
}

// test is a route group for testing purposes.
func (router *Router) test(r chi.Router) {
	// Use a logger middleware for logging requests
	r.Use(middleware.Logger)

	// Define a test endpoint for checking server health
	const prefix = "/api/v1/ping"
	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
}
