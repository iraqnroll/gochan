package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/views"
)

type Boards struct {
	Board  Template
	Thread Template

	PageData views.BasePageData

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
	b.PageData.BannerUrl, err = b.BoardService.GetBoardBannerUri(board.Uri)
	if err != nil {
		fmt.Printf("Failed to retrieve board banner : %s", err.Error())
	}

	fmt.Printf("Board banner path : %s", b.PageData.BannerUrl)

	b.PageData.PageData = board

	//If we get to this point we found a registered board, fill metadata and execute template
	b.Board.Execute(w, r, b.PageData)
}

func (b Boards) ThreadForm(w http.ResponseWriter, r *http.Request) {
	boardUri := chi.URLParam(r, "boardUri")
	threadId, err := strconv.Atoi(chi.URLParam(r, "threadId"))

	if err != nil {
		http.Error(w, "Invalid thread Id....", http.StatusBadRequest)
		return
	}

	result, err := b.BoardService.GetThread(threadId, boardUri)
	if err != nil {
		http.Error(w, "Unable to fetch thread...", http.StatusInternalServerError)
	}

	b.PageData.BannerUrl, err = b.BoardService.GetBoardBannerUri(boardUri)
	if err != nil {
		fmt.Printf("Failed to retrieve board banner : %s", err.Error())
	}

	fmt.Printf("Board banner path : %s", b.PageData.BannerUrl)

	b.PageData.PageData = result

	b.Thread.Execute(w, r, b.PageData)

}

func (b Boards) NewThread(w http.ResponseWriter, r *http.Request) {
	boardUri := chi.URLParam(r, "boardUri")
	board_id, err := strconv.Atoi(r.FormValue("boardId"))
	if err != nil {
		http.Error(w, "Invalid board id...", http.StatusBadRequest)
		fmt.Println("NewThread err: %w", err)
		return
	}

	_, err = b.BoardService.CheckBoard("", &board_id)
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

func (b Boards) NewReply(w http.ResponseWriter, r *http.Request) {

	thread_id, err := strconv.Atoi(chi.URLParam(r, "threadId"))
	if err != nil {
		http.Error(w, "Invalid thread_id entry, only numeric values are allowed.", http.StatusBadRequest)
		return
	}

	identifier := r.FormValue("replyIdentifier")
	content := r.FormValue("replyContent")

	err = b.BoardService.CreateReply(thread_id, identifier, content)
	if err != nil {
		http.Error(w, "Error while creating a reply....", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+chi.URLParam(r, "boardUri")+"/"+chi.URLParam(r, "threadId"), http.StatusFound)
}
