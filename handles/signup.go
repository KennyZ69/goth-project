package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"

	_ "github.com/lib/pq" // Import the PostgreSQL driver for it to work; ??? idk, it probably wont be used
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing the form: %v", err)
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		id, err := database.UsrId()
		role := r.FormValue("role")
		if err != nil {
			return fmt.Errorf("there was an error getting the user id: %v", err)
		}
		newUsr := database.User{
			Username: username,
			Email:    email,
			Id:       id,
			Details: database.UserProfileData{
				Role: role,
			},
		}

		exists, err := database.UserExists(newUsr)
		if err != nil {
			return fmt.Errorf("there was an error checking if the user exists: %v", err)
		}
		if exists {
			return fmt.Errorf("the user already exists")
		}

		if r.FormValue("password") != r.FormValue("confirm-password") {
			return fmt.Errorf("the password doesnt match")
		}
		newUsr.Password = database.HashPwd(r.FormValue("password"))
		newUsr.ComparePassword(r.FormValue("password"))
		if err = database.SaveUser(newUsr); err != nil {
			return fmt.Errorf("there was an error saving the user: %v", err)
		}

		tokenString, err := database.GenerateTokenString(newUsr)
		if err != nil {
			return fmt.Errorf("there was an error generating the token string: %v", err)
		}

		token, err := database.MakeToken(tokenString, newUsr.Id)
		if err != nil {
			return fmt.Errorf("there was an error making the token: %v", err)
		}

		if err = database.SaveToken(database.DB, token); err != nil {
			return fmt.Errorf("there was an error saving the token: %v", err)
		}

		// Set token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Expires:  token.ExpiresAt,
			HttpOnly: true,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = Render(components.AccountCreationSuccess(newUsr.Username), w, r); err != nil {
			return fmt.Errorf("there was an error rendering the success message: %v", err)
		}
		return nil

	} else if r.Method == "GET" {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("unsupported method")
}
