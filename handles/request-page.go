package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
	"time"
)

func HandleRequestPage(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
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
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})

		}

		err = HandleComponents(w, r)
		if err != nil {
			return fmt.Errorf("there was a problem handling the requests page: %v", err)
		}
		return nil
	}
	if r.Method == http.MethodPost {
		//TODO:Make the accept and deny button for the requests
		user, err := components.GetUserByCookie(r)
		if err != nil {
			return fmt.Errorf("problem getting user on post req on the req-page: %v", err)
		}
		requests, err := database.GetRequestsOnUser(user)
	}
	return fmt.Errorf("method not allowed on the requests page")
}
