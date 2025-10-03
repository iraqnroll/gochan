package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/models"
)

type Boards struct {
	Board Template

	BoardService *models.BoardService
}

func (b Boards) BoardForm(w http.ResponseWriter, r *http.Request) {
	boardUri := chi.URLParam(r, "boardUri")
	fmt.Println("potential board uri : ", boardUri)
	board, err := b.BoardService.GetBoard(boardUri)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No such board...", http.StatusNotFound)
			return
		}

		http.Error(w, "Unable to fetch requested board : "+err.Error(), http.StatusInternalServerError)
		return
	}

	//If we get to this point we found a registered board, fill metadata and execute template
	b.Board.Execute(w, r, board)

}
