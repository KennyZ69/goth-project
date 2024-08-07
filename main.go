package main

import (
	"gothstarter/handles"
	"gothstarter/middleware"
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

	router := chi.NewMux()

	handleHome := handles.MakeHandle(handles.HandleComponents)

	// router.Use(middleware.AuthMiddleware())

	router.Handle("/*", public())
	router.Get("/", middleware.AuthMiddleware(handleHome))
	router.Handle("/login", handles.MakeHandle(handles.HandleLogin))
	router.Handle("/signup", handles.MakeHandle(handles.HandleSignUp))
	router.Handle("/profile/{username}", handles.MakeHandle(handles.HandleProfile))
	router.Post("/logout", handles.MakeHandle(handles.HandleLogout))

	router.Handle("/finder", handles.MakeHandle(handles.HandleFinder))
	router.Handle("/finder/search", handles.MakeHandle(handles.HandleSearch))

	router.Handle("/api/countries", handles.MakeHandle(handles.CountryHandler))
	router.Handle("/api/cities", handles.MakeHandle(handles.CityHandler))

	router.Get("/team.html", handles.MakeHandle(handles.HandleTeamPage))

	listenAddr := os.Getenv("PORT")
	slog.Info("HTTP server was started", "listenAddr", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}
