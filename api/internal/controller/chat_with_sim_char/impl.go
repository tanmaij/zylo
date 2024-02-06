package chatwithsimchar

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

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

type GetCurrentConversation struct {
	ClientUUID string
}

func (i Impl) GetCurrentConversation(ctx context.Context, inp GetCurrentConversation) (model.Conversation, error) {
	c, err := i.conversationRedisRepo.Get(ctx, inp.ClientUUID)
	if err != nil {
		if errors.Is(err, conversationRedisRepo.ErrNotFound) {
			chars, err := i.characterRepo.List(ctx)
			if err != nil {
				return model.Conversation{}, err
			}

			charIdx := rand.Intn(len(chars))
			newConv := model.Conversation{
				Character: chars[charIdx],
				Messages:  []model.Message{},
			}

			if err := i.conversationRedisRepo.Set(ctx, inp.ClientUUID, time.Minute*3, newConv); err != nil {
				return model.Conversation{}, err
			}

			return newConv, nil
		}

		return model.Conversation{}, err
	}

	return c, nil
}

type ChatInput struct {
	ClientUUID string
	Message    string
}

type chatOutput struct {
	Sender       string `json:"sender"`
	SenderAvatar string `json:"senderAvatar"`
	Message      string `json:"message"`
}

func (i Impl) Chat(ctx context.Context, inp ChatInput) (model.Message, error) {
	c, err := i.conversationRedisRepo.Get(ctx, inp.ClientUUID)
	if err != nil {
		if errors.Is(err, conversationRedisRepo.ErrNotFound) {
			return model.Message{}, ErrConvNotFound
		}
		return model.Message{}, err
	}

	newMsg := model.Message{
		Role:    model.RoleUser,
		Content: inp.Message,
	}

	c.Messages = append(c.Messages, newMsg)
	if err := i.conversationRedisRepo.Set(ctx, inp.ClientUUID, time.Minute*3, c); err != nil {
		return model.Message{}, err
	}

	go func() {
		time.Sleep(time.Second * 4)

		i.Broadcast.EmitToClientUUID(inp.ClientUUID, ws.Message{
			EventName: "typing",
			Data:      "",
		})

		generated, err := i.chatCompletionClient.RequestToGenerate(context.Background(), openaiClient.ChatCompletionInput{
			Messages: append(append(c.Character.DefaultHistoryMsg, model.Message{
				Role:    model.RoleSystem,
				Content: c.Character.SystemMsg,
			}), c.Messages...),
		})
		if err != nil {
			i.Broadcast.EmitToClientUUID(inp.ClientUUID, ws.Message{
				EventName: "reset-typing",
				Data:      "",
			})
			return
		}

		for idx := range generated {
			marshal, err := json.Marshal(chatOutput{
				Sender:       c.Character.Name,
				SenderAvatar: c.Character.AvatarURL,
				Message:      generated[idx].Content,
			})
			if err != nil {
				continue
			}

			// Emit a message to the client with a success event and data.
			i.Broadcast.EmitToClientUUID(inp.ClientUUID, ws.Message{
				EventName: "chat",
				Data:      string(marshal),
			})
		}
	}()

	return newMsg, nil
}
