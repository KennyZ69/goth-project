package handles

import (
	"encoding/json"
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
	"time"
)

func HandleFinder(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		// Check if the auth_token cookie is present
		_, err := r.Cookie("auth_token")
		if err != nil {
			// If the cookie is not present, redirect to the login page
			http.Redirect(w, r, "/finder", http.StatusSeeOther)
			return nil
		}
		currentUser, err := components.GetUserByCookie(r)
		if err != nil {
			fmt.Printf("couldnt get the user on finder by cookie: %v", err)
		}
		token, err := database.IdToken(database.DB, currentUser.Id)
		if err != nil {
			return fmt.Errorf("there was an error getting the token by usr id in the finder: %v", err)
		}
		if token != "" {
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    token,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
				Path:     "/",
			})

		}
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("this method is not allowed")

}

func HandleSearch(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("q")
	currentUser, _ := components.GetUserByCookie(r)
	token, err := database.IdToken(database.DB, currentUser.Id)
	if err != nil {
		return fmt.Errorf("there was an error getting the token by usr id in the finder: %v", err)
	}
	if token != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
			Path:     "/",
		})

	}
	users, err := database.SearchUsers(database.DB, query, currentUser.Username)
	if err != nil {
		return fmt.Errorf("there was an error searching the users from database: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return nil
}
