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

	var results []database.Friend
	// querying right in the db to search the results
	// err = database.DB.QueryRow("SELECT * FROM friends WHERE user_id=$1 OR friend_id=$2 AND lower(friend_name) LIKE ?", currentUser.Id, currentUser.Id, "%"+query+"%").Scan(&results)
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM friends WHERE user_id=$1", currentUser.Id).Scan(&count)

	for i := 0; i < count; i++ {

		var result database.Friend
		err = database.DB.QueryRow("SELECT friend_id, friend_name FROM friends WHERE user_id=$1", currentUser.Id).Scan(&result.Friend_id, &result.Friend_name)
		if err == sql.ErrNoRows {
			return err
		}
		if err != nil {
			return err
		}
		fmt.Printf("friend name: %v;\n", result.Friend_name)
		results = append(results, result)
	}
	// if err == sql.ErrNoRows {
	// 	return fmt.Errorf("no rows in the database of friends for this user: %v", currentUser.Id)
	// } else if err != nil {
	// 	return fmt.Errorf("error somewhere in the db query: %v", err)
	// }

	var checkArrFriend []database.Friend
	w.Header().Set("Content-Type", "text/html")
	for _, f := range results {
		// 	// get the username of the friend
		// 	// err, friend := database.GetFriendById(f.Friend_id)
		// 	// if err != nil {
		// 	// return err
		// 	// }

		// if query == "" || strings.Contains(strings.ToLower(f.Friend_name), query) {
		if query != "" && strings.Contains(strings.ToLower(f.Friend_name), query) {
			fmt.Printf("found matching friend: %v\n", f.Friend_name)
			fmt.Fprintf(w, `
<li class="p-2 hover:bg-blue-100 cursor-pointer" hx-get="/api/openChatWithFriend?friendId=%v"
	hx-trigger="click" hx-target="this">
	%v	
</li>
		`, f.Friend_id, f.Friend_name)

			checkArrFriend = append(checkArrFriend, f)
		}
	}

	// Render a partial HTML template for the results
	// tmpl, err := template.ParseFiles("friend_search_results.html")
	// if err != nil {
	// 	return fmt.Errorf("error parsing the search results of friends: %v", err)
	// }
	// tmpl.Execute(w, results)
	if checkArrFriend == nil && query != "" {
		fmt.Fprintf(w, `
	<li class="p-2 text-gray-500">No results found</li>
	`)
	}

	return nil
}
