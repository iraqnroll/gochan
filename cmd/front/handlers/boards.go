package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
)

type Boards struct {
	BoardService   *services.BoardService
	FileService    *services.FileService
	ThreadsPerPage int
	ParentPage     models.ParentPageData
}

func NewBoardsHandler(boardSvc *services.BoardService, fileSvc *services.FileService, parentPage models.ParentPageData, threadsPerPage int) (b Boards) {
	b.BoardService = boardSvc
	b.FileService = fileSvc
	b.ThreadsPerPage = threadsPerPage
	b.ParentPage = parentPage

	return b
}

func (b Boards) Board(w http.ResponseWriter, r *http.Request) {
	board_uri := chi.URLParam(r, "board_uri")
	board, err := b.BoardService.GetBoard(board_uri)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Board not found...", http.StatusNotFound)
			return
		}

		http.Error(w, "Unable to fetch requested board : "+err.Error(), http.StatusInternalServerError)
		return
	}

	banner_uri, err := b.FileService.GetBoardBannerUri(board.Uri)
	if err != nil {
		fmt.Printf("Failed to retrieve board banner : %s", err.Error())
	}

	b.ParentPage.ChildViewModel = models.NewBoardViewModel(
		board.Id,
		b.ThreadsPerPage,
		board.Uri,
		board.Name,
		board.Description,
		banner_uri,
		board.Threads)

	views.Board(b.ParentPage).Render(r.Context(), w)
}
