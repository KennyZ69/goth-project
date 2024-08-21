package handles

import (
	"fmt"
	"gothstarter/database"
	"net/http"
	"time"
)

func HandleConnectRequest(w http.ResponseWriter, sender *database.User, receiver *database.User) error {

	var count int

	err := database.DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE sender_id=$1 AND receiver_id=$2 OR sender_id=$3 AND receiver_id=$4", sender.Id, receiver.Id, receiver.Id, sender.Id).Scan(&count)
	if err != nil {
		return fmt.Errorf("there was an error getting the number of these requests: %v", err)
	}

	if count > 0 {
		fmt.Printf("These users already have this request estabilished")
		return nil
	}

	_, err = database.DB.Exec("INSERT INTO connection_requests (sender_id, receiver_id, status, created_at) VALUES ($1, $2, $3, $4)", sender.Id, receiver.Id, database.StatusPending, time.Now())
	if err != nil {
		return fmt.Errorf("There was an error handling the connection request in the database: %v", err)
	}

	printResponse(w, receiver.Username, sender.Username)

	fmt.Printf("Connect request was successfully sent from %s to %s\n", sender.Username, receiver.Username)
	return nil
}

func printResponse(w http.ResponseWriter, receiverName string, senderName string) {
	fmt.Fprintf(w, `
   <div id="response-message"
	class="bg-green-100 flex flex-col justify-between items-center border border-green-400 text-green-700 gap-4 px-4 py-3 rounded relative mt-4 " role="alert">
	<strong class="font-bold">Success %s! </strong>
	<span class="block sm:inline">A connect request was sent to %s</span>
<button onclick="window.location.href='/profile/%s'"
           class="bg-transparent cursor-pointer border-0 text-black font-semibold outline-none focus:outline-none h-4 w-4 text-3xl self-center">
		×
        </button>	
</div>
    `, senderName, receiverName, receiverName)
}
