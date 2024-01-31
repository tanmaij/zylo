package ws

// BroadcastTo represents the target of the broadcast message.
type BroadcastTo string

// Constants representing different targets for broadcasting.
const (
	roomUUID   BroadcastTo = "room_uuid"
	clientUUID             = "client_uuid"
	clients                = "clients"

	// Emit except the sender.

	exRoomUUID = "ex_room_uuid"
	exClients  = "ex_clients"
)

// BroadcastMessage represents a message to be broadcast.
type BroadcastMessage struct {
	To         BroadcastTo // Target of the broadcast.
	ClientUUID string      // Client UUID for targeted messages.
	RoomUUID   string      // Room UUID for targeted room messages.
	Msg        Message     // The actual message to be broadcast.
}

// Broadcast is a channel for sending broadcast messages.
type Broadcast chan BroadcastMessage

// Emit sends a broadcast message to all connected clients.
func (b Broadcast) Emit(msg Message) {
	broadcastMsg := BroadcastMessage{To: clients, Msg: msg}
	b <- broadcastMsg
}

// EmitToClientUUID sends a broadcast message to a specific client by UUID.
func (b Broadcast) EmitToClientUUID(uuid string, msg Message) {
	broadcastMsg := BroadcastMessage{To: clientUUID, ClientUUID: uuid, Msg: msg}
	b <- broadcastMsg
}

// EmitToRoomUUID sends a broadcast message to all clients in a specific room by UUID.
func (b Broadcast) EmitToRoomUUID(uuid string, msg Message) {
	broadcastMsg := BroadcastMessage{To: roomUUID, RoomUUID: uuid, Msg: msg}
	b <- broadcastMsg
}
