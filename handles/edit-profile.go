package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"log"
	"net/http"
	"time"
)

func EditProfileHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("Invalid request method")
	}
	_, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		fmt.Printf("couldnt get the user on editing profile by cookie: %v", err)
	}
	token, err := database.IdToken(database.DB, currentUser.Id)
	if err != nil {
		return fmt.Errorf("there was an error getting the token by usr id to edit the profile: %v", err)
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

	err = r.ParseForm()
	if err != nil {
		return fmt.Errorf("Problem parsing the form data: %v", err)
	}

	for key, value := range r.Form {
		log.Printf("Form key: %s, value: %s", key, value)
	}
	newUsername := r.FormValue("username")
	country := r.FormValue("country-select")
	city := r.FormValue("city")
	role := r.FormValue("role")
	bio := r.FormValue("bio")

	if newUsername != currentUser.Username {
		_, err = database.DB.Exec(`
        UPDATE users 
        SET username = $1 
        WHERE user_id = $2`, newUsername, currentUser.Id)
		if err != nil {
			return fmt.Errorf("There was an error updating the username by his ID: %v", err)

		}
		currentUser.Username = newUsername
	}

	_, err = database.DB.Exec(`
	UPDATE user_details
	SET role = $1, bio = $2, country = $3, city = $4
	WHERE user_id = $5`, role, bio, country, city, currentUser.Id)
	if err != nil {
		return fmt.Errorf("There was an error updating the user_details: %v", err)
	}

	newUserData := database.UserProfileData{
		Bio:     bio,
		Role:    role,
		Country: country,
		City:    city,
	}
	currentUser.Details = newUserData

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Profile updated successfully for user: %s; %v\n", currentUser.Username, currentUser.Id)
	redirectUrl := fmt.Sprintf("/profile/%v", currentUser.Username)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}
