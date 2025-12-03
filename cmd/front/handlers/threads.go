package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
)

type Threads struct {
	ThreadService *services.ThreadService
	FileService   *services.FileService
	PostsPerPage  int
	ParentPage    models.ParentPageData
}

func NewThreadsHandler(threadSvc *services.ThreadService, fileSvc *services.FileService, parentPage models.ParentPageData, postsPerPage int) (t Threads) {
	t.ThreadService = threadSvc
	t.FileService = fileSvc
	t.ParentPage = parentPage
	t.PostsPerPage = postsPerPage

	return t
}

func (t Threads) Thread(w http.ResponseWriter, r *http.Request) {
	thread_id := chi.URLParam(r, "thread_id")
	board_uri := chi.URLParam(r, "board_uri")

	id, err := strconv.Atoi(thread_id)
	if err != nil {
		http.Error(w, "Invalid thread Id...", http.StatusBadRequest)
		return
	}
	thread, err := t.ThreadService.GetThread(id)
	if err != nil {
		http.Error(w, "Unable to fetch requested thread : "+err.Error(), http.StatusInternalServerError)
		return
	}
	banner_url, err := t.FileService.GetBoardBannerUri(board_uri)
	if err != nil {
		fmt.Printf("Failed to retrieve board banner : %s", err.Error())
	}

	t.ParentPage.ChildViewModel = models.NewThreadsViewModel(thread.Id, t.PostsPerPage, banner_url, board_uri, thread.Topic, thread.Posts[0], thread.Posts[1:])
	views.Thread(t.ParentPage).Render(r.Context(), w)
}
