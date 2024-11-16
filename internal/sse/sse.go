package sse

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

func HandleSSE(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Each connection registers its own message channel with the Broker's connection registry
	ev := make(chan event)

	generateRandomString := func() string {
		b := make([]byte, 10)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}

		return fmt.Sprintf("%x", b)
	}

	id := generateRandomString()
	// Signal the broker that we have a new connection in a specific room
	broker.newClients <- Client{
		ID:          id,
		MessageChan: ev,
	}

	// Disconnect the client when the connection is closed
	defer func() {
		// Signal the broker that the client has disconnected
		broker.closingClients <- Client{
			ID:          id,
			MessageChan: ev,
		}
	}()

	keepAliveTicker := time.NewTicker(15 * time.Second)
	defer keepAliveTicker.Stop()

	for {
		select {
		case <-r.Context().Done(): // Client disconnected
			return
		case message := <-ev: // Wait for a message
			fmt.Fprintf(w, "event: %s\n", message.Name)
			fmt.Fprintf(w, "data: %s\n\n", message.Data)

			// Send the data immediately
			flusher.Flush()
		case <-keepAliveTicker.C:
			fmt.Fprintf(w, "data:\n\n")
			flusher.Flush()

		}
	}
}
