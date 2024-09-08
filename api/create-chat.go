package api

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"
	"net/http"
)

func CreateChatHandle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method for creating chat")
	}
	fmt.Printf("handling the create-chat\n")
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		return err
	}
	friendId := r.URL.Query().Get("friendId")
	fmt.Printf("friendId = %v\n", friendId)
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM chats WHERE user1_id=$1 OR user1_id=$2 AND user2_id=$3 OR user2_id=$4", currentUser.Id, currentUser.Id, friendId, friendId).Scan(&count)
	if count < 1 {
		_, err = database.DB.Exec("INSERT INTO chats (user1_id, user2_id) VALUES ($1, $2)", currentUser.Id, friendId)
		fmt.Printf("new chat was created: %v and %v\n", currentUser.Id, friendId)
	}

	return nil
}
