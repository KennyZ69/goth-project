package main

import (
	"gothstarter/handles"
	"gothstarter/middleware"
	"gothstarter/ws"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	go ws.GlobalHub.Run()

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

	router.Use(middleware.AuthMiddleware)

	router.Handle("/*", public())
	router.Get("/", handleHome)
	// router.Get("/", middleware.AuthMiddleware(handleHome))
	router.Handle("/login", handles.MakeHandle(handles.HandleLogin))
	router.Handle("/signup", handles.MakeHandle(handles.HandleSignUp))
	router.Post("/logout", handles.MakeHandle(handles.HandleLogout))

	router.Handle("/profile/{username}", handles.MakeHandle(handles.HandleProfile))
	// router.Handle("/edit-profile", handles.MakeHandle(handles.EditProfileHandler))

	router.Handle("/finder", handles.MakeHandle(handles.HandleFinder))
	router.Handle("/finder/search", handles.MakeHandle(handles.HandleSearch))

	router.Handle("/requests/{username}", handles.MakeHandle(handles.HandleRequestPage))

	router.Handle("/ws", handles.MakeHandle(handles.WsHandler))
	router.Handle("/chat-try", handles.MakeHandle(handles.HandleChatTry))
	// router.Handle("/inbox/{username}", handles.MakeHandle(handles.HandleInbox))

	router.Handle("/api/countries", handles.MakeHandle(handles.CountryHandler))
	// router.Handle("/api/cities", handles.MakeHandle(handles.CityHandler))

	router.Get("/team.html", handles.MakeHandle(handles.HandleTeamPage))

	listenAddr := os.Getenv("PORT")
	slog.Info("HTTP server was started", "listenAddr", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}
