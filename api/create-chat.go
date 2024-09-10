package api

import (
	"fmt"
	"gothstarter/database"
	"gothstarter/handles"
	"gothstarter/layouts/components"
	"gothstarter/ws"
	"net/http"
	"strconv"
)

type ChatData struct {
	Messages    []*ws.Message  `json:"messages"`
	CurrentUser *database.User `json:"current_user"`
	Friend      *database.User `json:"friend"`
	ChatId      uint           `json:"chat_id"`
}

func CreateChatHandle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method for creating chat")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	fmt.Printf("handling the create-chat\n")
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		return err
	}

	friendIdString := r.URL.Query().Get("friendId")
	friendId, err := strconv.Atoi(friendIdString)
	fmt.Printf("friendId = %v\n", friendId)

	// var count int
	// err = tx.QueryRow("SELECT COUNT(*) FROM chats WHERE user1_id=$1 OR user1_id=$2 AND user2_id=$3 OR user2_id=$4", currentUser.Id, currentUser.Id, friendId, friendId).Scan(&count)

	var chatId uint
	err = tx.QueryRow("SELECT chat_id FROM chats WHERE (user1_id=$1 AND user2_id=$2) OR (user1_id=$2 AND user2_id=$1)", currentUser.Id, friendId).Scan(&chatId)

	// if count < 1 {
	if err != nil {

		err = tx.QueryRow("SELECT COALESCE(MAX(chat_id), 0) + 1 FROM chats").Scan(&chatId)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO chats (user1_id, user2_id, chat_id) VALUES ($1, $2, $3)", currentUser.Id, friendId, chatId)
		if err != nil {
			return fmt.Errorf("error inserting new chat into the db: %v", err)
		}
		fmt.Printf("new chat was created for: %v and %v\n", currentUser.Id, friendId)
	}

	//TODO: Get chat messages based on the chatId

	err = tx.QueryRow("SELECT chat_id FROM chats WHERE user1_id=$1 OR user2_id=$2 AND user1_id=$3 OR user2_id=$4", currentUser.Id, currentUser.Id, friendId, friendId).Scan(&chatId)
	if err != nil {
		return fmt.Errorf("error getting the chat id from chats: %v", err)
	}

	friend, err := database.GetUserById(uint(friendId))
	if err != nil {
		return err
	}

	// Fetch messages from the database
	var messages []*ws.Message
	rows, err := tx.Query("SELECT sender_id, sender_name, content, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at ASC", chatId)
	if err != nil {
		return fmt.Errorf("error fetching messages: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message ws.Message
		err := rows.Scan(&message.Client_id, &message.Username, &message.Text, &message.Created_at)
		if err != nil {
			return fmt.Errorf("error scanning message row: %v", err)
		}
		messages = append(messages, &message)
	}

	chatData := &ChatData{
		Messages:    messages,
		CurrentUser: currentUser,
		Friend:      friend,
		ChatId:      chatId,
	}

	handles.Render(Chat(*chatData), w, r)

	// tmpl, err := template.ParseFiles("api/html/chat-templ.html")
	// if err != nil {
	// 	return fmt.Errorf("error parsing the templ chat: %v", err)
	// }
	// var renderedChatTempl bytes.Buffer
	// err = tmpl.Execute(&renderedChatTempl, chatData)
	// if err != nil {
	// 	return fmt.Errorf("error rendering the chat templ into bytes buffer: %v\n", err)
	// }
	//
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// _, err = w.Write(renderedChatTempl.Bytes())
	// if err != nil {
	// 	return fmt.Errorf("error writing rendered template to response: %v", err)
	// }

	return tx.Commit()
}
