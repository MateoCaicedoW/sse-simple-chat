package sse

import (
	"fmt"
	"log"
	"time"
)

var msg = make(chan Event)
var broker = newBroker()

type Broker struct {
	rooms          map[string]map[string]chan Event // roomID -> clients
	newClients     chan Client
	closingClients chan string
	roomEntries    chan RoomEntry // Channel for entering a room
	roomLeaves     chan RoomLeave // Channel for leaving a room
	Notifier       chan Event
	clients        map[string]chan Event
}

type Client struct {
	ID          string
	MessageChan chan Event
}

type RoomEntry struct {
	ClientID string
	RoomID   string
}

type RoomLeave struct {
	ClientID string
	RoomID   string
}

func newBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan Event, 1),
		newClients:     make(chan Client),
		closingClients: make(chan string),
		roomEntries:    make(chan RoomEntry),
		roomLeaves:     make(chan RoomLeave),
		rooms:          make(map[string]map[string]chan Event),
		clients:        make(map[string]chan Event),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) listen() {
	for {
		select {
		case client := <-broker.newClients:
			// When a new client connects, we initialize their connection.
			log.Printf("Client %s connected", client.ID)
			broker.clients[client.ID] = client.MessageChan

		case clientID := <-broker.closingClients:
			// Remove the client from the clients map
			delete(broker.clients, clientID)
			log.Printf("Client %s disconnected", clientID)

			// When a client disconnects, remove them from all rooms.
			for roomID, clients := range broker.rooms {
				log.Printf("Client %s removed from room %s", clientID, roomID)
				// Clean up empty rooms
				if len(clients) == 0 {
					delete(broker.rooms, roomID)
				}
			}

		case entry := <-broker.roomEntries:
			// A client is joining a room
			if broker.rooms[entry.RoomID] == nil {
				broker.rooms[entry.RoomID] = make(map[string]chan Event)
			}
			broker.rooms[entry.RoomID][entry.ClientID] = broker.clients[entry.ClientID]
			log.Printf("Client %s entered room %s", entry.ClientID, entry.RoomID)

		case leave := <-broker.roomLeaves:
			// A client is leaving a room
			if clients, ok := broker.rooms[leave.RoomID]; ok {
				delete(clients, leave.ClientID)
				log.Printf("Client %s left room %s", leave.ClientID, leave.RoomID)
				if len(clients) == 0 {
					delete(broker.rooms, leave.RoomID)
				}
			}

		case event := <-broker.Notifier:
			// Broadcast event to all clients in the room
			if clients, ok := broker.rooms[event.RoomID]; ok {
				for clientID, clientMessageChan := range clients {
					fmt.Println("message", clientMessageChan)
					message, err := event.BuildMessage(map[string]interface{}{
						"content":       event.Content,
						"isCurrentUser": event.UserID == clientID,
					})
					if err != nil {
						log.Printf("Failed to build message for room %s: %v", event.RoomID, err)
						continue
					}
					event.Data = message
					select {
					case clientMessageChan <- event:
					case <-time.After(time.Second):
						log.Printf("Skipping slow/disconnected client in room %s", event.RoomID)
					}
				}
			}
		}
	}
}

// OpenSSEConnection opens a connection to the SSE server
func OpenSSEConnection() {
	go func() {
		for message := range msg {
			broker.Notifier <- message
		}
	}()
}

func RegisterClientToRoom(clientID, roomID string) {
	broker.roomEntries <- RoomEntry{ClientID: clientID, RoomID: roomID}
}

func UnregisterClientFromRoom(clientID, roomID string) {
	broker.roomLeaves <- RoomLeave{ClientID: clientID, RoomID: roomID}
}
