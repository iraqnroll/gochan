package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/models"
)

type Boards struct {
	Board  Template
	Thread Template

	BoardService *models.BoardService
}

func (b Boards) BoardForm(w http.ResponseWriter, r *http.Request) {
	boardUri := chi.URLParam(r, "boardUri")
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

func (b Boards) ThreadForm(w http.ResponseWriter, r *http.Request) {
	//boardUri := chi.URLParam(r, "boardUri")
	threadId, err := strconv.Atoi(chi.URLParam(r, "threadId"))

	if err != nil {
		http.Error(w, "Invalid thread Id....", http.StatusBadRequest)
		return
	}

	result, err := b.BoardService.GetThread(threadId)
	if err != nil {
		http.Error(w, "Unable to fetch thread...", http.StatusInternalServerError)
	}

	b.Thread.Execute(w, r, result)

}

func (b Boards) NewThread(w http.ResponseWriter, r *http.Request) {
	boardUri := chi.URLParam(r, "boardUri")
	board_id, err := b.BoardService.CheckBoard(boardUri)

	if err != nil {
		http.Error(w, "Error while resolving board...", http.StatusNotFound)
		fmt.Println("NewThread err: %w", err)
		return
	}

	//At this point we resolved a valid board, we can attempt to create a new thread.
	topic := r.FormValue("threadTopic")
	identifier := r.FormValue("threadIdentifier")
	content := r.FormValue("threadContent")

	_, err = b.BoardService.CreateThread(board_id, topic, identifier, content)

	if err != nil {
		http.Error(w, "Error while creating a new thread...", http.StatusInternalServerError)
		fmt.Println("NewThread err: %w", err)
		return
	}

	http.Redirect(w, r, "/"+boardUri, http.StatusFound)
}
