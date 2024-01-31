package ws

import (
	"log"

	chatwithsimchar "github.com/tanmaij/zylo/internal/controller/chat_with_sim_char"
	"github.com/tanmaij/zylo/pkg/ws"
)

type Handler struct {
	ChatwithsimcharController chatwithsimchar.Impl
}

func New(chatwithsimcharController chatwithsimchar.Impl) Handler {
	handler := Handler{ChatwithsimcharController: chatwithsimcharController}
	return handler
}

type pingInput struct {
	Message string `json:"message"`
}

func (h Handler) OnConnection(client ws.Client) error {
	log.Println(client.UUID, "connected")

	return nil
}

func (h Handler) OnDisconnection(client ws.Client) error {
	log.Println(client.UUID, "disconnected")

	return nil
}

func (h Handler) Ping(client ws.Client, msg string) error {
	h.ChatwithsimcharController.Test(client, chatwithsimchar.TestInput{
		Message: msg,
	})

	return nil
}
