package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"gothstarter/layouts/features"
	"net/http"
	"strings"
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
	query := r.URL.Query().Get("search")
	currentUser, _ := components.GetUserByCookie(r)
	users, err := database.SearchUsers(database.DB, query, currentUser.Username)
	if err != nil {
		return fmt.Errorf("there was an error searching the users from database: %v", err)
	}

	var results []database.User
	for _, usr := range users {
		if query != "" && contains(usr.Username, query) {
			results = append(results, usr)
		}
	}

	w.Header().Set("Content-Type", "text/html")
	Render(features.SearchResults(results), w, r)
	return nil
}

func contains(source, query string) bool {
	return strings.Contains(strings.ToLower(source), strings.ToLower(query))
}
