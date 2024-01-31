package ws

import "github.com/gorilla/websocket"

// Client represents a WebSocket client.
type Client struct {
	UUID      string          // Unique identifier for the client.
	RoomUUID  string          // Identifier for the room the client is in.
	Conn      *websocket.Conn // The WebSocket connection for the client.
	Broadcast Broadcast       // A channel for sending broadcast messages to the client.
}
