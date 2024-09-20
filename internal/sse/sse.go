package sse

import (
	"fmt"
	"net/http"
	"simple-chat-sse/internal/auth"

	"github.com/leapkit/leapkit/core/server/session"
)

func HandleSSE(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Get the room ID from the parameters
	roomID := r.PathValue("id")

	// Each connection registers its own message channel with the Broker's connection registry
	ev := make(chan Event)

	// Get the client ID (authenticated user)
	clientID := auth.GetUserID(session.FromCtx(r.Context()))

	// Signal the broker that we have a new connection in a specific room
	broker.newClients <- Client{
		ID:          clientID,
		MessageChan: ev,
		RoomID:      roomID,
	}

	// Disconnect the client when the connection is closed
	defer func() {
		// Signal the broker that the client has disconnected
		broker.closingClients <- Client{
			ID:          clientID,
			RoomID:      roomID,
			MessageChan: ev,
		}
	}()

	for {
		select {
		case <-r.Context().Done(): // Client disconnected
			return
		case message := <-ev: // Wait for a message
			fmt.Fprintf(w, "event: %s\n", message.Name)
			fmt.Fprintf(w, "data: %s\n\n", message.Data)

			// Send the data immediately
			flusher.Flush()
		}
	}
}
