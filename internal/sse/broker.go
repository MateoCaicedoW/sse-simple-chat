package sse

import (
	"log"
	"time"
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
	messages chan Event
}

func newBroker() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		messages:       make(chan Event),
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
		case client := <-b.newClients:
			b.clients[client.ID] = client
		case client := <-b.closingClients:
			delete(b.clients, client.ID)
		case message := <-b.messages:
			for _, client := range b.clients {
				select {
				case client.MessageChan <- message:
				case <-time.After(time.Second * 5): // Timeout for slow clients
					delete(b.clients, client.ID)
					close(client.MessageChan)
				default:
					log.Printf("Message dropped for client: %s", client.ID)
				}
			}
		}
	}
}
