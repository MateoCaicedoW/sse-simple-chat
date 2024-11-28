package sse

import (
	"fmt"
	"strings"
)

type event struct {
	// Name is the event name
	Name string

	// Data is the event data
	Data string
}

// NewEvent creates a new event
func NewEvent(name string) *event {
	return &event{
		Name: name,
	}
}

// BuildMessage builds the message
func (e event) BuildMessage(data interface{}) (string, error) {
	message, err := RenderToString("message.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to render message: %v", err)
	}

	return strings.Replace(message, "\n", "", -1), nil
}
