package handlers

import (
	"net/http"

	"github.com/iraqnroll/gochan/db/models"
	"github.com/iraqnroll/gochan/db/services"
	"github.com/iraqnroll/gochan/views"
)

type Home struct {
	BoardService        *services.BoardService
	PostService         *services.PostService
	ParentPage          models.ParentPageData
	NumberOfRecentPosts int
}

func NewHomeHandler(boardService *services.BoardService, postService *services.PostService, parentPage models.ParentPageData, numOfRecentPosts int) (h Home) {
	h.BoardService = boardService
	h.PostService = postService
	h.ParentPage = parentPage
	h.NumberOfRecentPosts = numOfRecentPosts
	return h
}

func (h Home) Home(w http.ResponseWriter, r *http.Request) {
	boards, err := h.BoardService.GetBoardList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	recentPosts, err := h.PostService.GetMostRecent(h.NumberOfRecentPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.ParentPage.ChildViewModel = models.NewHomeViewModel(boards, recentPosts, nil)

	views.Home(h.ParentPage).Render(r.Context(), w)
}
