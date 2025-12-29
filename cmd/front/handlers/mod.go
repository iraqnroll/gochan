package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
)

const (
	DEFAULT_USERS_REDIRECT = "/mod/users"
)

type Mod struct {
	UserService   *services.UserService
	BoardService  *services.BoardService
	ThreadService *services.ThreadService
	FileService   *services.FileService
	ParentPage    models.ParentPageData
}

func NewModHandler(
	userSvc *services.UserService,
	boardSvc *services.BoardService,
	threadSvc *services.ThreadService,
	fileSvc *services.FileService,
	parentPage models.ParentPageData) (m Mod) {
	m.UserService = userSvc
	m.BoardService = boardSvc
	m.ThreadService = threadSvc
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

func (m Mod) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
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
