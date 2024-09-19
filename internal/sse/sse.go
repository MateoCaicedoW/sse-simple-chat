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
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Each connection registers its own message channel with the Broker's connections registry
	ev := make(chan Event)

	// Signal the broker that we have a new connection
	broker.newClients <- Client{ID: auth.GetUserID(session.FromCtx(r.Context())), MessageChan: ev}

	// Listen to connection close and un-register messageChan

	defer func() {
		// Signal the broker that we are done
		broker.closingClients <- auth.GetUserID(session.FromCtx(r.Context()))
	}()

	for {
		select {
		case <-r.Context().Done(): // Client has disconnected
			return
		case message := <-ev: // Wait for a message
			fmt.Fprintf(w, "event: %s\n", message.Name)
			fmt.Fprintf(w, "data: %s\n\n", message.Data)

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()
		}
	}
}
