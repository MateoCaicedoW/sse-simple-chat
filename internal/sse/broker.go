package sse

import (
	"log"
)

var broker = newBroker()

type Broker struct {
	// Map of clients connected to the broker
	clients map[string]Client

	// Register a new client
	newClients chan Client

	// Unregister a client
	closingClients chan Client

	// Send a message to all clients
	messages chan event
}

type Client struct {
	// ID is the client id
	ID string

	// MessageChan is the message channel for the client to receive messages
	MessageChan chan event
}

func newBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		messages:       make(chan event),
		newClients:     make(chan Client),
		closingClients: make(chan Client),
		clients:        map[string]Client{},
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return
}

// Listen for new clients and messages
func (b *Broker) listen() {
	for {
		select {
		// When we receive a new client, we add it to the map
		case client := <-b.newClients:
			b.clients[client.ID] = client

		// When a client has gone, we delete it from the map
		case client := <-b.closingClients:
			delete(b.clients, client.ID)

		// When we receive a new event, we send it to all clients
		case message := <-b.messages:
			for _, client := range b.clients {
				select {
				case client.MessageChan <- message:
				default:
					log.Printf("Message dropped for client: %s", client.ID)
				}
			}
		}
	}
}
