package sse

import (
	"log"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range clients` started.
const patience time.Duration = time.Second * 1

var msg = make(chan Event)
var broker = newBroker()

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan Event

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan string

	// Client connections registry
	clients map[string]chan Event
}

type Client struct {
	ID          string
	MessageChan chan Event
	RoomID      string
}

type Room struct {
	ID          string
	MessageChan chan Event
}

func newBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan Event, 1),
		newClients:     make(chan Client),
		closingClients: make(chan string),
		clients:        make(map[string]chan Event),
	}

	// Set it running - listening and broadcasting events
	go listen(broker)

	return
}

func listen(broker *Broker) {
	for {
		select {
		case s := <-broker.newClients:
			// A new client has connected.
			//create a new client identifier
			clientID := s.ID
			broker.clients[clientID] = s.MessageChan
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// Broadcast to all clients
			for id, clientMessageChan := range broker.clients {
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

func OpenSSEConnection() {
	go func() {
		for message := range msg {
			broker.Notifier <- message
		}
	}()
}
