package sse

import (
	"fmt"
	"strings"
)

type Event struct {
	// Name is the event name
	Name string

	// Data is the event data
	Data string

	//Content is the content of the message
	Content string

	// Client is the client that the event is sent to
	ClientID string

	//UserID is the user id
	UserID string
}

// NewEvent creates a new event
func NewEvent(name string) *Event {
	return &Event{
		Name:     name,
		Data:     "",
		ClientID: "",
		UserID:   "",
	}
}

// BuildMessage builds the message
func (e Event) BuildMessage(data interface{}) (string, error) {
	message, err := RenderToString("message.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to render message: %v", err)
	}

	return strings.Replace(message, "\n", "", -1), nil
}

func (e *Event) Broadcast() {
	msg <- *e
}
