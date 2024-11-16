package messages

import (
	"cmp"
	"fmt"
	"net/http"
	"os"

	"github.com/leapkit/leapkit/core/render"
)

var baseURL = cmp.Or(os.Getenv("BASE_URL"), "http://localhost:3000")

func LoadAudio(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	rw.Set("audio", fmt.Sprintf("%s/attachments/%s", baseURL, r.PathValue("audio")))
	err := rw.Render("messages/load_audio.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
