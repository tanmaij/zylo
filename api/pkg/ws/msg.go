package ws

// Message represents a WebSocket message.
type Message struct {
	EventName string `json:"event_name"` // The name of the WebSocket event.
	Data      string `json:"data"`       // The data associated with the WebSocket event.
}
