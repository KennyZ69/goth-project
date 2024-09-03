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
			err = saveFriendInDb(uint(sender_int), currentUser.Id)
			if err != nil {
				return err
			}

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<div class="bg-blue-500 text-white px-4 py-2 rounded">Accepted</div>
			<button onclick="refreshPage()" type="button"
					class="p-1 ml-auto bg-transparent cursor-pointer border-0 text-black float-right text-3xl leading-none font-semibold outline-none focus:outline-none">
					<span
						class="bg-transparent text-black h-6 w-6 text-2xl block outline-none focus:outline-none">
						×
					</span>
				</button>
			<script>
				function refreshPage(){
				window.location.href="%s"
				}
			</script>
				`, r.URL.Path)

			return nil
		case "deny":
			// Handle deny logic
			err := database.UpdateReqStatus(uint(sender_int), currentUser.Id, database.StatusDenied)
			if err != nil {
				return fmt.Errorf("there was a problem updating the request status: %v", err)
			}

			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<div class="bg-blue-500 text-white px-4 py-2 rounded">Denied</div>
			<button onclick="refreshPage()" type="button"
					class="p-1 ml-auto bg-transparent cursor-pointer border-0 text-black float-right text-3xl leading-none font-semibold outline-none focus:outline-none">
					<span
						class="bg-transparent text-black h-6 w-6 text-2xl block outline-none focus:outline-none">
						×
					</span>
				</button>
			<script>
				function refreshPage(){
				window.location.href="%s"
				}
			</script>
				`, r.URL.Path)

			return nil
		default:
			return fmt.Errorf("somehow invalid operation with buttons")
		}
	}
	return fmt.Errorf("method not allowed on the requests page")
}

func saveFriendInDb(user1_id, user2_id uint) error {
	_, err := database.DB.Exec("INSERT INTO friends (user_id, friend_id, created_at) VALUES ($1, $2, $3)", user1_id, user2_id, time.Now())
	if err != nil {
		return fmt.Errorf("there was an error saving the friend to database: %v", err)
	}
	return nil
}
