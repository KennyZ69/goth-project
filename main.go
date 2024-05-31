package main

import (
	"gothstarter/handles"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	router.Handle("/*", public())
	router.Get("/", handles.MakeHandle(handles.HandleHome))

	// http.Handle("/*", public())
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	handles.Render(home.Index("home"), w, r)
	// })
	// http.Handle("/home", handles.MakeHandle(handles.HandleHome))

	listenAddr := os.Getenv("PORT")
	slog.Info("HTTP server was started", "listenAddr", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
	// log.Fatal(http.ListenAndServe(listenAddr, nil))
}
