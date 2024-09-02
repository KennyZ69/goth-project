package handles

import (
	"fmt"
	"net/http"
)

func HandleChatTry(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		err := HandleComponents(w, r)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("method not allowed on chat-try")
}

func HandleInbox(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		err := HandleComponents(w, r)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("not supported method yet on inbox")
}
