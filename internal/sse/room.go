package sse

type Room struct {
	// ID is the room id
	ID string

	// Clients are the clients in the room
	clients map[string]chan Event
}

// Add client to the room
func (r *Room) addClient(client Client) {
	r.clients[client.ID] = client.MessageChan
}

// Remove client from the room
func (r *Room) removeClient(clientID string) {
	delete(r.clients, clientID)
}

type Client struct {
	// ID is the client id
	ID string

	// MessageChan is the message channel for the client to receive messages
	MessageChan chan Event

	// RoomID is the room id the client is in
	RoomID string
}
