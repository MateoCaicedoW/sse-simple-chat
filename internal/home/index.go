package home

import (
	"fmt"
	"net/http"
	"net/url"
	"simple-chat-sse/internal/auth"
	"simple-chat-sse/internal/sse"

	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server/session"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	session := session.FromCtx(r.Context())
	auth.SetUserID(session)

	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		chatID = "1"
	}

	last, err := url.Parse(r.Referer())
	if err != nil {
		fmt.Println(err)
	}

	if last.Query().Get("chatID") != "" {
		sse.RegisterClientToRoom(auth.GetUserID(session), last.Query().Get("chatID"))
	}

	rw.Set("chatID", chatID)
	sse.RegisterClientToRoom(auth.GetUserID(session), chatID)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Push-URL", r.URL.String())
		err := rw.RenderClean("home/chat.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = rw.Render("home/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
