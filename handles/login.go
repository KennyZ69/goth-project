package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"

	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("error parsing the form: %v", err)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		usr, err := database.GetUserByName(database.DB, username)
		if err != nil {
			return fmt.Errorf("there was an error getting the user: %v", err)
		}
		if err = usr.ComparePassword(password); err != nil {
			return fmt.Errorf("invalid password")
		}
		tokenString, err := database.GenerateTokenString(*usr)
		if err != nil {
			return fmt.Errorf("failed to generate token string; err: %v", err)
		}

		if err = database.DeleteUserTokens(database.DB, usr.Id); err != nil {
			return fmt.Errorf("failed to delete user tokens")
		}

		sesToken, err := database.MakeToken(tokenString, usr.Id)
		if err != nil {
			return fmt.Errorf("there was an error making the token: %v", err)
		}
		if err = database.SaveToken(database.DB, sesToken); err != nil {
			return fmt.Errorf("there was an error saving the token: %v", err)
		}

		// Set token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Expires:  sesToken.ExpiresAt,
			HttpOnly: true,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = Render(components.LoginSuccess(usr.Username), w, r); err != nil {
			return fmt.Errorf("there was an error rendering the success message: %v", err)
		}

		return nil
	} else if r.Method == http.MethodGet {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("this method is not allowed")
}
