package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		err := HandleComponents(w, r)
		if err != nil {

			return fmt.Errorf("problem handling the profile: %v", err)
		}
		return nil
	}
	if r.Method == http.MethodPost {
		nameFromPath := r.URL.Path[len("/profile/"):]
		userFromPath, err := database.GetUserByName(database.DB, nameFromPath)
		if err != nil {
			return err
		}
		// userPathDetails, err := database.GetDetailsById(database.DB, userFromPath.Id)
		// if err != nil {
		// 	return err
		// }
		// userFromPath.Details = *userPathDetails

		currentUser, err := components.GetUserByCookie(r)
		if currentUser.Id == userFromPath.Id {
			err = EditProfileHandler(w, r)
			if err != nil {
				return fmt.Errorf("problem handling edit-profile on the profile handler: %v", err)
			}
			return nil
		} else {
			//TODO: now the user is checking ElseProfile so there would be the connect handle
			sender := currentUser
			receiver := userFromPath
			err = HandleConnectRequest(w, sender, receiver)
			if err != nil {
				return err
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	return fmt.Errorf("this method is not allowed")
}
