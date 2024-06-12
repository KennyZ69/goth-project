package handles

import (
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		// reqUsr := database.User{
		// 	Username: r.FormValue("username"),
		// 	Email:    r.FormValue("email"),
		// 	Password: database.HashPwd(r.FormValue("password")),
		// }

		return nil
	} else if r.Method == http.MethodGet {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("this method is not allowed")
}
