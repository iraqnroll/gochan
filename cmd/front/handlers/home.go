package handlers

import (
	"net/http"

	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
)

type Home struct {
	BoardService *services.BoardService
	ParentPage   models.ParentPageData
}

func NewHomeHandler(boardService *services.BoardService, parentPage models.ParentPageData) (h Home) {
	h.BoardService = boardService
	h.ParentPage = parentPage
	return h
}

func (h Home) Home(w http.ResponseWriter, r *http.Request) {
	boards, err := h.BoardService.GetBoardList()
	if err != nil {
		http.Error(w, "Failed to fetch boards.", http.StatusInternalServerError)
	}
	h.ParentPage.ChildViewModel = models.NewHomeViewModel(boards, nil)

	views.Home(h.ParentPage).Render(r.Context(), w)
}
