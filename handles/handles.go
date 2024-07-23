package handles

import (
	"fmt"
	"gothstarter/layouts/components"
	index "gothstarter/layouts/index"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type httpHandler func(w http.ResponseWriter, r *http.Request) error

func Render(t templ.Component, w http.ResponseWriter, r *http.Request) error {
	return t.Render(r.Context(), w)
}

func MakeHandle(h httpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("http handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func HandleComponents(w http.ResponseWriter, r *http.Request) error {
	isAuthenticated := components.GetAuth(r)
	if isAuthenticated {
		user, err := components.GetUserByCookie(r)
		if err != nil {
			fmt.Printf("There was an error handling the GetUserByCookie request: %v", err)
		}
		//TODO: Make handle logout function to be able to logout using the button from the dropdown
		return Render(index.Index(r, isAuthenticated, user.Username), w, r)
	}
	return Render(index.Index(r, isAuthenticated, ""), w, r)
}

func Static(filename string, w http.ResponseWriter, r *http.Request) error {
	// Parse the template file
	tmpl, err := template.ParseFiles("html/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Render the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func HandleTeamPage(w http.ResponseWriter, r *http.Request) error {
	err := Static("team.html", w, r)
	if err != nil {
		return fmt.Errorf("there was an error with handling the team page")
	}
	return nil
}
