package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
	"strconv"
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
		currentUser, err := components.GetUserByCookie(r)
		if err != nil {
			return fmt.Errorf("problem getting the user on post on the req-page: %v", err)
		}
		r.ParseForm()

		action := r.FormValue("action")
		fmt.Printf("action: %v\n", action)
		sender_id := r.FormValue("sender_id")
		fmt.Printf("sender_id: %v\n", sender_id)
		sender_int, err := strconv.Atoi(sender_id)
		if err != nil {
			return fmt.Errorf("problem parsing the sender_id to int: %v", err)
		}

		switch action {
		case "accept":
			// Handle accept logic
			err := database.UpdateReqStatus(uint(sender_int), currentUser.Id, database.StatusAccepted)
			if err != nil {
				return fmt.Errorf("there was a problem updating the request status: %v", err)
			}
			fmt.Fprintf(w, "Request from sender %s accepted.", sender_id)
			return nil
		case "deny":
			// Handle deny logic
			err := database.UpdateReqStatus(uint(sender_int), currentUser.Id, database.StatusDenied)
			if err != nil {
				return fmt.Errorf("there was a problem updating the request status: %v", err)
			}
			fmt.Fprintf(w, "Request from sender %s denied.", sender_id)
			return nil
		default:
			return fmt.Errorf("somehow invalid operation with buttons")
		}
	}
	return fmt.Errorf("method not allowed on the requests page")
}
