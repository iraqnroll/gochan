package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/db/services"
)

type API struct {
	UserService *services.UserService
}

// List godoc
//
//	@summary        List users
//	@description    List all registered users
//	@tags           users
//	@accept         json
//	@produce        json
//	@success        200 {array}     []models.User
//	@failure        500 {object}    err.Error
//	@router         /user/all [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	result, err := a.UserService.GetAll()
	if err != nil {
		http.Error(w, "Error while retrieving user list", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error while converting response to JSON", http.StatusInternalServerError)
		return
	}
}
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user id specified.", http.StatusBadRequest)
		return
	}
	result, err := a.UserService.GetUserById(user_id)
	if err != nil {
		http.Error(w, "Error while retrieving user.", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error while converting response to JSON", http.StatusInternalServerError)
		return
	}
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {}
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {}
func (a *API) Update(w http.ResponseWriter, r *http.Request) {}
func (a *API) Login(w http.ResponseWriter, r *http.Request)  {}
func (a *API) Logout(w http.ResponseWriter, r *http.Request) {}
