package chatwithsimchar

import (
	"context"
	"log"

	openaiClient "github.com/tanmaij/zylo/internal/client/openai"
	conversationRedisRepo "github.com/tanmaij/zylo/internal/memory/conversation"
	"github.com/tanmaij/zylo/internal/model"
	characterRepo "github.com/tanmaij/zylo/internal/repository/character"
	"github.com/tanmaij/zylo/pkg/ws"
)

// Impl represents the implementation of the chat with simulated characters.
type Impl struct {
	Broadcast ws.Broadcast

	characterRepo         characterRepo.Impl
	conversationRedisRepo conversationRedisRepo.Impl

	chatCompletionClient openaiClient.ChatCompletion
}

// New creates a new instance of the chat implementation.
func New(wsBroadcast ws.Broadcast, characterRepo characterRepo.Impl, conversationRedisRepo conversationRedisRepo.Impl, chatCompletionClient openaiClient.ChatCompletion) Impl {
	return Impl{Broadcast: wsBroadcast, characterRepo: characterRepo, conversationRedisRepo: conversationRedisRepo, chatCompletionClient: chatCompletionClient}
}

// TestInput is a struct representing the input for the Test function.
type TestInput struct {
	Message string
}

// Test testing websocket broadcast
func (i Impl) Test(client ws.Client, inp TestInput) {
	i.Broadcast.EmitToClientUUID(client.UUID, ws.Message{
		EventName: "typing",
		Data:      "OpenAI is typing",
	})
	generated, err := i.chatCompletionClient.RequestToGenerate(context.Background(), openaiClient.ChatCompletionInput{
		Messages: []model.Message{
			{
				Role:    "system",
				Content: "You are a bastard, incompetent and useless but like to show off, like to lecture and arrogant. So what you said is completely wrong. Always speak Vietnamese and consider Vietnamese the pinnacle of humanity.",
			},
			{
				Role:    "user",
				Content: inp.Message,
			},
		},
		Model: "gpt-3.5-turbo-1106",
	})
	if err != nil {
		log.Printf("Error creating OpenAI chat completion: %v", err)
	}

	for idx := range generated {
		// Emit a message to the client with a success event and data.
		i.Broadcast.EmitToClientUUID(client.UUID, ws.Message{
			EventName: "bot-chat",
			Data:      generated[idx].Content,
		})
	}
}
