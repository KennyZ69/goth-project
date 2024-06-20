package handles

import "net/http"

func HandleProfile(w http.ResponseWriter, r *http.Request) error {
	err := HandleComponents(w, r)
	if err != nil {
		return err
	}
	return nil
}
