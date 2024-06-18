package components

import "net/http"

type ComponentsData struct {
	isAuthenticated bool `json:"isAuthenticated"`
}

func GetAuth(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		// Handle the case where isAuthenticated is not set
		isAuthenticated = false
	}
	// return ComponentsData{
	// 	isAuthenticated: isAuthenticated,
	// }
	return isAuthenticated
}
