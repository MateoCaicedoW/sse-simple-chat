package sse

type Client struct {
	// ID is the client id
	ID string

	// MessageChan is the message channel for the client to receive messages
	MessageChan chan Event
}
