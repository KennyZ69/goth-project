package handles

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
	"time"
)

func EditProfileHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("error parsing the form in edit-profile: %v", err)
	}
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
		return fmt.Errorf("there was an error getting the token by user id to edit the profile: %v", err)
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

	// for key, value := range r.Form {
	// 	log.Printf("Form key: %s, value: %s", key, value)
	// }
	newUsername := r.FormValue("username")
	country := r.FormValue("country-select")
	city := r.FormValue("city")
	role := r.FormValue("role")
	bio := r.FormValue("bio")

	var newUser database.User
	newUser = *currentUser

	newUserData := database.UserProfileData{
		Bio:     bio,
		Role:    role,
		Country: country,
		City:    city,
	}
	newUser.Details = newUserData
	newUser.Username = newUsername

	if err := UpdateUser(newUser, currentUser.Username); err != nil {
		return fmt.Errorf("There was an error updating the user profile in the database: %v", err)
	}
	// Return a success message as an HTML snippet
	fmt.Fprintf(w, `
   <div id="response-message"
	class="bg-green-100 flex flex-col justify-between items-center border border-green-400 text-green-700 gap-4 px-4 py-3 rounded relative mt-4 " role="alert">
	<strong class="font-bold">Success!  </strong>
	<span class="block sm:inline">Your profile has been updated %s.</span>
<button onclick="window.location.href='/profile/%s'"
           class="bg-transparent cursor-pointer border-0 text-black font-semibold outline-none focus:outline-none h-4 w-4 text-3xl self-center">
		×
        </button>	
</div>
    `, newUsername, newUsername)
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Profile updated successfully for user: %s; %v\n", currentUser.Username, currentUser.Id)
	return nil
}

func UpdateUser(user database.User, oldName string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if user.Username != oldName {
		if valid := database.UsernameValidator(user.Username); !valid {
			return fmt.Errorf("there was an error validating the username")
		}

		if exists := database.UsernameExists(user.Username); exists {
			return fmt.Errorf("this username already exists")
		}
		_, err = tx.Exec("UPDATE users SET username = $1 WHERE user_id = $2",
			user.Username, user.Id)
		if err != nil {
			return err
		}

	}

	// Update the user_details table
	_, err = tx.Exec("UPDATE user_details SET bio = $1, profile_image = $2, role = $3, country = $4, city = $5 WHERE user_id = $6",
		user.Details.Bio, user.Details.ProfileImage, user.Details.Role, user.Details.Country, user.Details.City, user.Id)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()

}
