package handles

import (
	"fmt"
	"gothstarter/ws"
	"net/http"
)

func WsHandler(w http.ResponseWriter, r *http.Request) error {
	err := ws.SocketHandler(w, r)
	if err != nil {
		fmt.Printf("error on the WsHandler(): %v\n", err)
	}
	return nil
}
