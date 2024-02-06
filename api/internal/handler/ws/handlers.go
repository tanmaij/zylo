package ws

import (
	"log"

	chatwithsimchar "github.com/tanmaij/zylo/internal/controller/chat_with_sim_char"
	"github.com/tanmaij/zylo/pkg/utils"
	"github.com/tanmaij/zylo/pkg/ws"
)

// Handler is responsible for handling WebSocket events.
type Handler struct {
	ChatwithsimcharController chatwithsimchar.Impl
}

// New creates a new instance of the WebSocket handler.
func New(chatwithsimcharController chatwithsimchar.Impl) Handler {
	handler := Handler{ChatwithsimcharController: chatwithsimcharController}
	return handler
}

// pingInput is a struct representing the input for the Ping function.
type pingInput struct {
	Message string `json:"message"`
}

type connectedOutput struct {
	UUID string `json:"uuid"`
}

// OnConnection is called when a client establishes a WebSocket connection.
func (h Handler) OnConnection(client ws.Client) error {
	var output = connectedOutput{
		UUID: client.UUID,
	}
	rs, err := utils.AnyToJSON(output)
	if err != nil {
		return err
	}

	client.Broadcast.EmitToClientUUID(client.UUID, ws.Message{EventName: "connection", Data: string(rs)})

	return nil
}

// OnDisconnection is called when a client disconnects from the WebSocket.
func (h Handler) OnDisconnection(client ws.Client) error {
	log.Println(client.UUID, "disconnected")

	return nil
}

// Ping handles the "ping" event from the client, simulating a chat interaction with a simulated character.
func (h Handler) Ping(client ws.Client, msg string) error {
	// Call the Test function from ChatwithsimcharController to simulate a chat interaction.
	h.ChatwithsimcharController.Test(client, chatwithsimchar.TestInput{
		Message: msg,
	})

	return nil
}
