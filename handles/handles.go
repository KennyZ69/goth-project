package handles

import (
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
	return Render(index.Index(r), w, r)
}
