package chatwithsimchar

import (
	conversationRedisRepo "github.com/tanmaij/zylo/internal/memory/conversation"
	characterRepo "github.com/tanmaij/zylo/internal/repository/character"
	"github.com/tanmaij/zylo/pkg/ws"
)

// Impl represents the implementation of the chat with simulated characters.
type Impl struct {
	Broadcast             ws.Broadcast
	characterRepo         characterRepo.Impl
	conversationRedisRepo conversationRedisRepo.Impl
}

// New creates a new instance of the chat implementation.
func New(wsBroadcast ws.Broadcast, characterRepo characterRepo.Impl, conversationRedisRepo conversationRedisRepo.Impl) Impl {
	return Impl{wsBroadcast, characterRepo, conversationRedisRepo}
}

// TestInput is a struct representing the input for the Test function.
type TestInput struct {
	Message string
}

// Test testing websocket broadcast
func (i Impl) Test(client ws.Client, inp TestInput) {
	// Emit a message to the client with a success event and data.
	i.Broadcast.EmitToClientUUID(client.UUID, ws.Message{
		EventName: "successfully",
		Data:      "You're the best",
	})
}
