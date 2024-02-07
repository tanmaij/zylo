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

type typingOutput struct {
	Name string `json:"name"`
}

func (i Impl) Chat(ctx context.Context, inp ChatInput) (model.Message, error) {
	ca, err := i.conversationRedisRepo.Get(ctx, inp.ClientUUID)
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

	ca.Messages = append(ca.Messages, newMsg)
	if err := i.conversationRedisRepo.Overwrite(ctx, inp.ClientUUID, ca); err != nil {
		return model.Message{}, err
	}

	go func(currentConv model.Conversation) {
		time.Sleep(time.Second * 4)

		marshal, err := json.Marshal(typingOutput{Name: currentConv.Character.Name})
		if err != nil {
		}

		i.Broadcast.EmitToClientUUID(inp.ClientUUID, ws.Message{
			EventName: "typing",
			Data:      string(marshal),
		})

		generated, err := i.chatCompletionClient.RequestToGenerate(context.Background(), openaiClient.ChatCompletionInput{
			Messages: append(append(currentConv.Character.DefaultHistoryMsg, model.Message{
				Role:    model.RoleSystem,
				Content: currentConv.Character.SystemMsg,
			}), currentConv.Messages...),
		})

		rsMarshal, err := json.Marshal(typingOutput{Name: currentConv.Character.Name})
		if err != nil {

		}
		i.Broadcast.EmitToClientUUID(inp.ClientUUID, ws.Message{
			EventName: "reset-typing",
			Data:      string(rsMarshal),
		})

		if err != nil {
			return
		}

		for idx := range generated {
			currentConv.Messages = append(currentConv.Messages, generated[idx])
			if err := i.conversationRedisRepo.Overwrite(context.Background(), inp.ClientUUID, currentConv); err != nil {
				log.Printf("overwriting got err: %v", err)
			}

			marshal, err := json.Marshal(chatOutput{
				Sender:       currentConv.Character.Name,
				SenderAvatar: currentConv.Character.AvatarURL,
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
	}(ca)

	return newMsg, nil
}
