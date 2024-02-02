package main

import (
	"log"
	"net/http"

	"github.com/tanmaij/zylo/cmd/serverd/router"
	"github.com/tanmaij/zylo/config"
	"github.com/tanmaij/zylo/internal/client/openai"
	chatwithsimchar "github.com/tanmaij/zylo/internal/controller/chat_with_sim_char"
	wsHandler "github.com/tanmaij/zylo/internal/handler/ws"
	conversationRedisRepo "github.com/tanmaij/zylo/internal/memory/conversation"
	characterRepo "github.com/tanmaij/zylo/internal/repository/character"
	"github.com/tanmaij/zylo/pkg/redis"
	"github.com/tanmaij/zylo/pkg/ws"

	"github.com/go-chi/chi/v5"
)

// main function serves as the entry point for the application.
func main() {
	// Initialize application configuration.
	if err := config.Initial(); err != nil {
		log.Fatal("Failed to initialize configuration:", err)
	}

	// Retrieve application configuration instance.
	appCfg := config.Instance.App

	// Initialize a new Redis client.
	_, err := redis.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to initialize Redis client:", err)
	}

	// Initialize WebSocket listener with CORS check function.
	websocketListener := ws.NewWebSocketListener(func(r *http.Request) bool { return true })
	wsBroadcast := websocketListener.Broadcast

	// Start a goroutine to handle incoming WebSocket messages.
	go websocketListener.HandleMessages()

	// Initialize repositories
	charRepo := characterRepo.New()
	converRedisRepo := conversationRedisRepo.New()

	// Initialize api clients
	chatCompletation, err := openai.NewChatCompletion(config.Instance.OpenAI.APIKey)
	if err != nil {
		log.Fatalf("Error creating chat completion, error: %v", err)
	}

	// Create a new controller for handling chat with simulated characters.
	chatWithSimCharCtrl := chatwithsimchar.New(wsBroadcast, charRepo, converRedisRepo, chatCompletation)

	// Create a WebSocket handler with the chat controller.
	wsHandler := wsHandler.New(chatWithSimCharCtrl)

	// Create a new Chi router.
	r := chi.NewRouter()

	// Create a new router and configure routes.
	appRouter := router.New(websocketListener, wsHandler)
	appRouter.RegisterRoutes(r)

	// Start the HTTP server and listen for incoming requests.
	http.ListenAndServe(appCfg.Port, r)
}
