package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/jabuta/dreampicai/handler"
	"github.com/jabuta/dreampicai/pkg/db"
	"github.com/jabuta/dreampicai/pkg/sb"

	"github.com/joho/godotenv"
)

//go:embed public
var fs embed.FS

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {

	if err := initEverything(embedMigrations); err != nil {
		log.Fatal(err)
	}
	defer db.Conf.PgxConneciton.Close(db.Conf.Ctx)

	router := chi.NewMux()
	router.Use(handler.WithUser)
	// router.Get("/path",handler)
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(fs))))
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
		auth.Get("/account", handler.MakeHandler(handler.HandleAccountIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)

	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything(migrations embed.FS) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.InitDatabase(migrations); err != nil {
		return err
	}
	if err := sb.Init(); err != nil {
		return err
	}
	return nil
}
