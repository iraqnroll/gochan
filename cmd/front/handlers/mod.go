package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/iraqnroll/gochan/db/models"
	"github.com/iraqnroll/gochan/db/services"
	"github.com/iraqnroll/gochan/views"
)

const (
	DEFAULT_USERS_REDIRECT  = "/mod/users"
	DEFAULT_BOARDS_REDIRECT = "/mod/boards"
)

type Mod struct {
	UserService   *services.UserService
	BoardService  *services.BoardService
	ThreadService *services.ThreadService
	PostService   *services.PostService
	FileService   *services.FileService
	ParentPage    models.ParentPageData
}

func NewModHandler(
	userSvc *services.UserService,
	boardSvc *services.BoardService,
	threadSvc *services.ThreadService,
	postSvc *services.PostService,
	fileSvc *services.FileService,
	parentPage models.ParentPageData) (m Mod) {
	m.UserService = userSvc
	m.BoardService = boardSvc
	m.ThreadService = threadSvc
	m.PostService = postSvc
	m.FileService = fileSvc
	m.ParentPage = parentPage

	return m
}

func (m Mod) ModPage(w http.ResponseWriter, r *http.Request) {
	views.ModPageComponent(m.ParentPage).Render(r.Context(), w)
}

func (m Mod) ModUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.UserService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.ParentPage.ChildViewModel = models.NewModUsersViewModel(users)
	views.ModUsersPageComponent(m.ParentPage).Render(r.Context(), w)
}

func (m Mod) ModBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := m.BoardService.GetAdminBoardList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.ParentPage.ChildViewModel = models.NewModBoardsViewModel(boards)
	views.ModBoardsPageComponent(m.ParentPage).Render(r.Context(), w)
}

func (m Mod) EditBoardPage(w http.ResponseWriter, r *http.Request) {
	board_id, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board, err := m.BoardService.GetBoardById(board_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.ParentPage.ChildViewModel = models.NewModBoardViewModel(*board)
	views.ModUpdateBoardPageComponent(m.ParentPage).Render(r.Context(), w)
}

func (m Mod) EditUserPage(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := m.UserService.GetUserById(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.ParentPage.ChildViewModel = models.NewModUserViewModel(*user)
	views.ModUpdateUserPageComponent(m.ParentPage).Render(r.Context(), w)
}

func (m Mod) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.User
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = m.UserService.CreateNew(model.Username, model.Password, model.Email, model.User_type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, DEFAULT_USERS_REDIRECT, http.StatusFound)
}

func (m Mod) CreateBoard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.Board
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = m.BoardService.Create(model.Uri, model.Name, model.Description, model.OwnerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, DEFAULT_BOARDS_REDIRECT, http.StatusFound)
}

func (m Mod) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.UserService.Delete(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, DEFAULT_USERS_REDIRECT, http.StatusFound)
}

func (m Mod) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	board_id, err := strconv.Atoi(r.FormValue("boardId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	board_uri := r.FormValue("boardUri")

	err = m.BoardService.Delete(board_id, board_uri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, DEFAULT_BOARDS_REDIRECT, http.StatusFound)
}

// TODO: Handle static content folder changes relative to board metadata changes (only URI update necessary for now.)
func (m Mod) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	board_id, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.BoardDto
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = m.BoardService.UpdateBoard(board_id, model.Uri, model.Name, model.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, DEFAULT_BOARDS_REDIRECT, http.StatusFound)

}

func (m Mod) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var model models.User
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(&model, r.PostForm)
	if err != nil {
		fmt.Printf("Failed to decode form : %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashed_pass, err := m.UserService.HashPassword(model.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	model.Password_hash = string(hashed_pass)

	_, err = m.UserService.UpdateUser(user_id, model.User_type, model.Username, model.Password_hash, model.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, DEFAULT_USERS_REDIRECT, http.StatusFound)
}

func (m Mod) SoftDeletePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.PostService.SoftDeletePost(post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (m Mod) RemoveSoftDeleteFromPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.PostService.RemoveSoftDeleteFromPost(post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (m Mod) PinThread(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread_id, err := strconv.Atoi(r.FormValue("thread_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.ThreadService.PinThread(thread_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (m Mod) UnpinThread(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread_id, err := strconv.Atoi(r.FormValue("thread_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.ThreadService.UnpinThread(thread_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (m Mod) PinPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.PostService.PinPost(post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (m Mod) UnpinPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := r.FormValue("redirect")

	err = m.PostService.UnpinPost(post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}
