package handles

import (
	"fmt"
	"gothstarter/database"
	"log"
	"net/http"
	"time"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) error {
	// Retrieve the auth token from the cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Print("Problem getting the cookie")
		return fmt.Errorf("could not retrieve auth_token cookie: %v", err)
	}
	tokenString := cookie.Value

	// Delete the token from the database
	if err = database.DeleteTokenByToken(database.DB, tokenString); err != nil {
		log.Print("Problem deleting the token from the user_tokens db")
		return fmt.Errorf("failed to delete the token from the database: %v", err)
	}

	// Invalidate the cookie by setting its expiration time to the past
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Redirect to the home page
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
