package sse

import (
	"log"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range clients` started.
const patience time.Duration = time.Second * 1

var broker = newBroker()

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	notifier chan Event

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan Client

	// Rooms registry for clients
	rooms map[string]*Room
}

// Get or create a room
// If the room does not exist, it will be created and returned
// If the room exists, it will be returned
func (b *Broker) getOrCreateRoom(roomID string) *Room {
	room, exists := b.rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			clients: make(map[string]chan Event),
		}
		b.rooms[roomID] = room
	}
	return room
}

func newBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		notifier:       make(chan Event, 1),
		newClients:     make(chan Client),
		closingClients: make(chan Client),
		rooms:          make(map[string]*Room),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			// A new client has connected.
			room := broker.getOrCreateRoom(s.RoomID)
			room.addClient(s)

			log.Printf("Client added. %d registered clients", len(room.clients))
		case s := <-broker.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			if room, exists := broker.rooms[s.RoomID]; exists {
				room.removeClient(s.ID)
				log.Printf("Removed client from room. %d registered clients", len(room.clients))
			}

			if len(broker.rooms[s.RoomID].clients) == 0 {
				delete(broker.rooms, s.RoomID)
				log.Printf("Removed room. %d registered rooms", len(broker.rooms))
			}

		case event := <-broker.notifier:
			room, exists := broker.rooms[event.RoomID]
			if !exists {
				log.Printf("Room does not exist")
				continue
			}

			// Broadcast to all clients
			for id, clientMessageChan := range room.clients {
				message, err := event.BuildMessage(map[string]interface{}{
					"content":       event.Content,
					"isCurrentUser": event.UserID == id,
				})

				if err != nil {
					log.Printf("Failed to build message: %v", err)
				}

				event.Data = message
				select {
				case clientMessageChan <- event:
				case <-time.After(patience):
					log.Printf("Skipping client with ID: %s (slow or disconnected)", id)
				}
			}
		}
	}
}
