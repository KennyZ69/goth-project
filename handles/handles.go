package handles

import (
	"fmt"
	"gothstarter/layouts/components"
	index "gothstarter/layouts/index"
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
	user, err := components.GetUserByCookie(w, r)
	if err != nil {
		fmt.Printf("There was an error handling the GetUserByCookie request: %v", err)
	}

	data := components.GetAuth(r)
	fmt.Println(data)
	return Render(index.Index(r, data, user.Username), w, r)
}
