package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jabuta/dreampicai/handler"
	"github.com/jabuta/dreampicai/pkg/sb"

	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)
	// router.Get("/path",handler)
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/login", handler.MakeHandler(handler.HandleLogInIndex))
	router.Get("/login/provider/google", handler.MakeHandler(handler.HandleLogInWithGoogle))
	router.Post("/login", handler.MakeHandler(handler.HandleLogInCreate))
	router.Get("/signup", handler.MakeHandler(handler.HandleSignupIndex))
	router.Post("/signup", handler.MakeHandler(handler.HandleSignupCreate))
	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))
	router.Get("/auth/callback/pkce", handler.MakeHandler(handler.HandleAuthCallbackPKCE))
	router.Post("/log-out", handler.MakeHandler(handler.HandleLogoutCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth)
		auth.Get("/settings", handler.MakeHandler(handler.HandleSettingsIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)

	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return sb.Init()
}
