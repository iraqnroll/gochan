package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/db/models"
	"github.com/iraqnroll/gochan/db/services"
	"github.com/iraqnroll/gochan/views"
)

type Threads struct {
	ThreadService *services.ThreadService
	PostService   *services.PostService
	FileService   *services.FileService
	PostsPerPage  int
	ParentPage    models.ParentPageData
}

func NewThreadsHandler(threadSvc *services.ThreadService, postSvc *services.PostService, fileSvc *services.FileService, parentPage models.ParentPageData, postsPerPage int) (t Threads) {
	t.ThreadService = threadSvc
	t.PostService = postSvc
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

	//TODO: Only set for_mod for specific user types.
	//Fetch user from context
	for_mod := context.User(r.Context()) != nil

	thread, err := t.ThreadService.GetThread(id, for_mod)
	if err != nil {
		http.Error(w, "Unable to fetch requested thread : "+err.Error(), http.StatusInternalServerError)
		return
	}
	banner_url, err := t.FileService.GetBoardBannerUri(board_uri)
	if err != nil {
		fmt.Printf("Failed to retrieve board banner : %s", err.Error())
	}

	model := models.NewThreadsViewModel(thread.Id, t.PostsPerPage, banner_url, board_uri, thread.Topic, thread.Posts[0], thread.Posts[1:], false, thread.Pinned)
	t.renderThreadMarkdown(model)

	t.ParentPage.ChildViewModel = model

	views.Thread(t.ParentPage).Render(r.Context(), w)
}

func (t Threads) Reply(w http.ResponseWriter, r *http.Request) {
	board_uri := chi.URLParam(r, "board_uri")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.PostDto
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Add validation before saving the new reply
	//Validate attached media formats.
	m := r.MultipartForm
	files := m.File["file-input"]

	if len(files) > 0 {
		err, result := t.FileService.CheckForInvalidFileFormats(files)
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
	model.Post_fprint = t.PostService.GenerateFingerprint(GetClientIp(r))

	//Split tripcode password from poster name and generate hash if exists.
	model.Identifier, model.Tripcode = t.PostService.GetTripcodeHash(model.Identifier, context.User(r.Context()) != nil)

	new_post, err := t.PostService.CreatePost(model.ThreadId, model.Identifier, model.Content, model.Post_fprint, model.Tripcode, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Handle attached media.
	attached_media, err := t.FileService.HandleFileUploads(files, board_uri, new_post.ThreadId, new_post.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	og_media := t.FileService.GetFilenames(files)

	err = t.PostService.UpdateAttachedMedia(new_post.Id, attached_media, og_media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/%s/%d", board_uri, new_post.ThreadId), http.StatusFound)
}

func (t Threads) renderThreadMarkdown(m models.ThreadViewModel) {
	m.OPPost.Content, _ = t.PostService.RenderSafeMarkdown(m.OPPost.Content)

	for i := range m.Replies {
		m.Replies[i].Content, _ = t.PostService.RenderSafeMarkdown(m.Replies[i].Content)
	}
}
