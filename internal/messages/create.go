package messages

import (
	"net/http"
	"simple-chat-sse/internal/auth"
	"simple-chat-sse/internal/sse"
	"strings"

	"github.com/leapkit/leapkit/core/server/session"
)

func Create(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")

	msg = strings.TrimSpace(msg)
	msg = strings.Replace(msg, "\n", "<br>", -1)

	message := sse.NewEvent("message")
	message.UserID = auth.GetUserID(session.FromCtx(r.Context()))
	message.Content = msg

	message.Broadcast()
}
