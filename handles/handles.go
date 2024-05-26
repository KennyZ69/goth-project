package handles

import (
	"net/http"

	"github.com/a-h/templ"
)

func HandleTempl(t templ.Component, w http.ResponseWriter, r *http.Request) {
	t.Render(r.Context(), w)
}
