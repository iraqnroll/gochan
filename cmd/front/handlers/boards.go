package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
	"github.com/microcosm-cc/bluemonday"
)

type Boards struct {
	BoardService   *services.BoardService
	ThreadService  *services.ThreadService
	FileService    *services.FileService
	ThreadsPerPage int
	ParentPage     models.ParentPageData
	PostPolicy     *bluemonday.Policy
}

func NewBoardsHandler(boardSvc *services.BoardService, threadSvc *services.ThreadService, fileSvc *services.FileService, parentPage models.ParentPageData, threadsPerPage int, pPol *bluemonday.Policy) (b Boards) {
	b.BoardService = boardSvc
	b.ThreadService = threadSvc
	b.FileService = fileSvc
	b.ThreadsPerPage = threadsPerPage
	b.ParentPage = parentPage
	b.PostPolicy = pPol

	return b
}

func (b Boards) Board(w http.ResponseWriter, r *http.Request) {
	board_uri := chi.URLParam(r, "board_uri")
	for_mod := context.User(r.Context()) != nil

	board, err := b.BoardService.GetBoard(board_uri, for_mod)
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
		fmt.Printf("Failed to retrieve board banner : %s \n", err.Error())
	}

	b.ParentPage.ChildViewModel = models.NewBoardViewModel(
		board.Id,
		b.ThreadsPerPage,
		board.Uri,
		board.Name,
		board.Description,
		banner_uri,
		board.Threads,
		b.PostPolicy)

	views.Board(b.ParentPage).Render(r.Context(), w)
}

func (b Boards) NewThread(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.ThreadDto
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Implement validation before saving anything
	//Validate attached media formats.
	m := r.MultipartForm
	files := m.File["file-input"]

	if len(files) > 0 {
		err, result := b.FileService.CheckForInvalidFileFormats(files)
		if !result {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, "Invalid media formats in attachments.", http.StatusBadRequest)
			return
		}
	}

	//Generate fingerprint
	model.Posts[0].Post_fprint = b.ThreadService.GenerateFingerprint(GetClientIp(r))

	new_thread, err := b.ThreadService.CreateThread(
		model.BoardId,
		model.Topic,
		model.Posts[0].Identifier,
		model.Posts[0].Content,
		model.Posts[0].Post_fprint)
	if err != nil {
		fmt.Printf("Failed to create a new thread : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Handle attached media
	attached_media, err := b.FileService.HandleFileUploads(files, model.BoardUri, new_thread.Id, new_thread.Posts[0].Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	og_media := b.FileService.GetFilenames(files)

	err = b.ThreadService.UpdateAttachedMedia(new_thread.Posts[0].Id, attached_media, og_media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%d", model.BoardUri, new_thread.Id), http.StatusFound)
}
