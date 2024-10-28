package home

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

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
