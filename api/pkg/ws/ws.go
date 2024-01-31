package ws

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebsocketListener handles WebSocket connections and messages.
type WebsocketListener struct {
	Clients   []Client
	Broadcast Broadcast

	upgrader             *websocket.Upgrader
	eventHandlers        map[string]func(client Client, data string) error
	connectionHandler    func(client Client) error
	disconnectionHandler func(client Client) error
}

// NewWebSocketListener creates a new WebsocketListener with the specified CORS check function.
func NewWebSocketListener(cors func(r *http.Request) bool) *WebsocketListener {
	return &WebsocketListener{
		upgrader:      &websocket.Upgrader{CheckOrigin: cors},
		Clients:       []Client{},
		eventHandlers: make(map[string]func(client Client, data string) error),
		Broadcast:     make(Broadcast),
	}
}

// On registers an event handler for a specific event name.
func (l *WebsocketListener) On(eventName string, handler func(client Client, data string) error) {
	l.eventHandlers[eventName] = handler
}

// OnConnection sets the handler function for new connections.
func (l *WebsocketListener) OnConnection(handler func(client Client) error) {
	l.connectionHandler = handler
}

// OnDisconnection sets the handler function for disconnections.
func (l *WebsocketListener) OnDisconnection(handler func(client Client) error) {
	l.disconnectionHandler = handler
}

// HandleMessages is a routine to send messages to clients.
func (l *WebsocketListener) HandleMessages() {
	for {
		b := <-l.Broadcast

		log.Printf("Received message: %v", b)

		switch b.To {
		case clientUUID:
			for i := range l.Clients {
				if l.Clients[i].UUID == b.ClientUUID {
					l.Clients[i].Conn.WriteJSON(b.Msg)
				}
			}
		case roomUUID:
			for i := range l.Clients {
				if l.Clients[i].RoomUUID == b.RoomUUID {
					l.Clients[i].Conn.WriteJSON(b.Msg)
				}
			}
		case clients:
			for i := range l.Clients {
				l.Clients[i].Conn.WriteJSON(b.Msg)
			}
		default:
		}
	}
}

// register creates a new client and registers it if it's not already registered.
func (l *WebsocketListener) register(cnn *websocket.Conn) Client {
	for i := range l.Clients {
		if l.Clients[i].Conn == cnn {
			return l.Clients[i]
		}
	}

	uuid := uuid.New().String()
	client := Client{
		UUID:      uuid,
		RoomUUID:  uuid,
		Conn:      cnn,
		Broadcast: l.Broadcast,
	}

	l.Clients = append(l.Clients, client)
	l.connectionHandler(client)

	return client
}

// removeClientByConn removes a client from the list based on its connection.
func (l *WebsocketListener) removeClientByConn(cnn *websocket.Conn) {
	newSlice := []Client{}
	for i := range l.Clients {
		if l.Clients[i].Conn == cnn {
			continue
		}

		newSlice = append(newSlice, l.Clients[i])
	}

	l.Clients = newSlice
}

// HandleConnections handles WebSocket connections and incoming messages.
func (l *WebsocketListener) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := l.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register client
	client := l.register(ws)

	for {
		var msg Message
		if err := client.Conn.ReadJSON(&msg); err != nil {
			log.Printf("Cannot read the message: %v", err)
			l.disconnectionHandler(client)
			l.removeClientByConn(ws)
			break
		}

		if msg.EventName == "" {
			log.Printf("Event name not found")
			continue
		}

		handler, ok := l.eventHandlers[msg.EventName]
		if !ok {
			log.Printf("No event matched: %v", msg.EventName)
			continue
		}

		if err := handler(client, msg.Data); err != nil {
			log.Printf("Handling event %v got an error: %v", msg.EventName, err)
		}
	}
}
