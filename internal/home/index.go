package home

import (
	"net/http"
	"simple-chat-sse/internal/auth"

	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server/session"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	auth.SetUserID(session.FromCtx(r.Context()))

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		chatID = "1"
	}

	rw.Set("chatID", chatID)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Push-URL", r.URL.String())
		err := rw.RenderClean("home/chat.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := rw.Render("home/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
