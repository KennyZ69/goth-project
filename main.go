package main

import (
	"gothstarter/handles"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Run the PostgreSQL setup script
	cmd := exec.Command("/bin/bash", "./setup_postgres.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute setup_postgres.sh: %v", err)
	}

	// // if err := goose.Up(database.DB, "./script.sql"); err != nil {
	// // 	log.Fatal(err)
	// // }

	router := chi.NewMux()

	router.Handle("/*", public())
	router.Get("/", handles.MakeHandle(handles.HandleComponents))
	router.Handle("/login", handles.MakeHandle(handles.HandleLogin))
	router.Handle("/signup", handles.MakeHandle(handles.HandleSignUp))

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
