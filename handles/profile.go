package handles

import (
	"fmt"
	"net/http"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		err := HandleComponents(w, r)
		if err != nil {

			return fmt.Errorf("problem handling the profile: %v", err)
		}
		return nil
	}
	return fmt.Errorf("this method is not allowed")
}
