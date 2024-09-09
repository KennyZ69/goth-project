package api

import (
	"database/sql"
	"fmt"
	"gothstarter/database"
	"gothstarter/layouts/components"

	"net/http"
	"strings"
)

// Here will be the functions for api usage for the inbox/chat

func SearchFriends(w http.ResponseWriter, r *http.Request) error {
	currentUser, err := components.GetUserByCookie(r)
	if err != nil {
		return err
	}
	query := r.URL.Query().Get("friendSearch")
	query = strings.ToLower(query)
	fmt.Printf("query: %v\n", query)

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	var results []database.Friend

	rows, err := tx.Query("SELECT friend_id, friend_name FROM friends WHERE user_id=$1", currentUser.Id)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var result database.Friend
		err := rows.Scan(&result.Friend_id, &result.Friend_name)
		if err != nil {
			return fmt.Errorf("error scanning rows for friend: %v", err)
		}

		fmt.Printf("friend name: %v;\n", result.Friend_name)
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "text/html")
	for _, f := range results {

		if query != "" && strings.Contains(strings.ToLower(f.Friend_name), query) {
			fmt.Printf("found matching friend: %v\n", f.Friend_name)
			fmt.Fprintf(w, `
<li class="p-2 hover:bg-blue-100 cursor-pointer" hx-get="/api/openChat?friendId=%v"
	hx-trigger="click" hx-target="#chatContent" hx-swap="innerHTML">
	%v	
</li>
		`, f.Friend_id, f.Friend_name)

		}
	}

	if results == nil && query != "" {
		fmt.Fprintf(w, `
	<li class="p-2 text-gray-500">No results found</li>
	`)
	}

	return nil
}
