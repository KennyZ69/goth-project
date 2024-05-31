package main

import (
	"gothstarter/handles"
	"gothstarter/layouts/components/home"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	listenAddr := os.Getenv("PORT")

	router := chi.NewMux()

	router.Handle("/*", public())

	// router.Handle("/", handles.MakeHandle(handles.HandleHome))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handles.Render(home.Index("home"), w, r)

	})

	// http.Handle("/*", public())
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	handles.Render(home.Index("home"), w, r)
	// })
	// http.Handle("/home", handles.MakeHandle(handles.HandleHome))

	slog.Info("HTTP server was started", "listenAddr", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
