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

	chats := []Chat{
		{ID: "1", Name: "Friends", Pic: "/public/group1.webp"},
		{ID: "2", Name: "Family", Pic: "/public/family.jpeg"},
	}

	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		chatID = chats[0].ID
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = chats[0].Name
	}

	rw.Set("chatID", chatID)
	rw.Set("name", name)
	rw.Set("chats", chats)

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

type Chat struct {
	ID   string
	Name string
	Pic  string
}
