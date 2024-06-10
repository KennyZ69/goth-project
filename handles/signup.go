package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"

	_ "github.com/lib/pq" // Import the PostgreSQL driver for it to work; ??? idk, it probably wont be used
)

// TODO : connect to db, add users etc...
// var db = database.InitDB()

func HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	email := r.FormValue("email")
	if r.Method == "POST" {
		newUsr := database.User{
			Username: username,
			Email:    email,
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
		err = database.SaveUser(newUsr)
		if err != nil {
			return fmt.Errorf("there was an error saving the user: %v", err)
		}
		if err = Render(components.AccountCreationSuccess(), w, r); err != nil {
			return fmt.Errorf("there was an error rendering the success message: %v", err)
		}
		return nil
	} else if r.Method == "GET" {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("unsupported method")
}
