package handlers

import (
	"net/http"

	"github.com/iraqnroll/gochan/views"
)

type Home struct {
	Footer *views.FooterData
}

func (h Home) Home(w http.ResponseWriter, r *http.Request) {
	views.Home(*h.Footer).Render(r.Context(), w)
}
