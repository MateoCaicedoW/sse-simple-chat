package messages

import (
	"cmp"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"simple-chat-sse/internal/sse"
	"strings"
	"time"

	"github.com/leapkit/leapkit/core/server"
)

var uploadsFolder = cmp.Or(os.Getenv("UPLOADS_FOLDER"), "./uploads")

func Create(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")
	if msg == "" {
		return
	}
	msg = strings.TrimSpace(msg)
	msg = strings.Replace(msg, "\n", "<br>", -1)

	event := sse.NewEvent("chat")
	content, err := event.BuildMessage(map[string]interface{}{
		"content": msg,
	})

	if err != nil {
		server.Error(w, err, http.StatusInternalServerError)
		return
	}

	event.Data = content
	sse.Broadcast(event)
}

func Audio(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Could not read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%d.wav", time.Now().Unix())

	os.MkdirAll(uploadsFolder, os.ModePerm) // Ensure the upload directory exists
	dst, err := os.Create(filepath.Join(uploadsFolder, fileName))
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	event := sse.NewEvent("chat")
	event.Data = fmt.Sprintf(`<div hx-get="/load-audio/%s" hx-trigger="load" hx-swap="outerHTML"></div>`, fileName)
	sse.Broadcast(event)
}
