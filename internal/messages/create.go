package messages

import (
	"net/http"
	"simple-chat-sse/internal/sse"
	"strings"

	"github.com/leapkit/leapkit/core/server"
)

func Create(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")

	msg = strings.TrimSpace(msg)
	msg = strings.Replace(msg, "\n", "<br>", -1)

	message := sse.NewEvent("chat")
	content, err := message.BuildMessage(map[string]interface{}{
		"content": msg,
	})

	if err != nil {
		server.Error(w, err, http.StatusInternalServerError)
		return
	}

	message.Data = content
	message.Broadcast()
}
