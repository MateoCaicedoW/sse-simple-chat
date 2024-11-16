package internal

import (
	"cmp"
	"embed"
	"net/http"
	"os"

	"simple-chat-sse/internal/home"
	"simple-chat-sse/internal/messages"
	"simple-chat-sse/internal/sse"

	"simple-chat-sse/public"

	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
)

var (
	//go:embed **/*.html **/*.html *.html
	tmpls embed.FS

	// DB is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DB = db.ConnectionFn(
		cmp.Or(os.Getenv("DATABASE_URL"), "database.db?_timeout=5000&_sync=1"),
		db.WithDriver("sqlite3"),
	)
)

// Server interface exposes the methods
// needed to start the server in the cmd/app package
type Server interface {
	Addr() string
	Handler() http.Handler
}

func New() Server {
	// Creating a new server instance with the
	// default host and port values.
	r := server.New(
		server.WithHost(cmp.Or(os.Getenv("HOST"), "0.0.0.0")),
		server.WithPort(cmp.Or(os.Getenv("PORT"), "3000")),
		server.WithSession(
			cmp.Or(os.Getenv("SESSION_SECRET"), "d720c059-9664-4980-8169-1158e167ae57"),
			cmp.Or(os.Getenv("SESSION_NAME"), "leapkit_session"),
		),
		server.WithAssets(public.Files),
	)

	r.Use(render.Middleware(
		render.TemplateFS(tmpls, "internal"),

		render.WithDefaultLayout("layout.html"),
	))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("GET /{$}", home.Index)
	r.HandleFunc("GET /chat", sse.HandleSSE)
	r.HandleFunc("POST /message", messages.Create)
	r.HandleFunc("POST /audio", messages.Audio)
	r.HandleFunc("GET /load-audio/{audio}", messages.LoadAudio)
	r.Folder("/attachments", os.DirFS("./uploads"))
	return r
}
