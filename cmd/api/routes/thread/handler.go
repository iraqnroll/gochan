package thread

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/services"
)

type API struct {
	ThreadService *services.ThreadService
}

// List godoc
//
//	@summary        Get thread
//	@description    Get a specific thread by id
//	@tags           threads
//	@accept         json
//	@produce        json
//	@success        200 {object}     models.ThreadDto
//	@failure        500 {object}    err.Error
//	@router         /threads/{id} [get]
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid thread id", http.StatusBadRequest)
		return
	}

	result, err := a.ThreadService.GetThread(id, false)
	if err != nil {
		http.Error(w, "Failed to fetch thread by id", http.StatusInternalServerError)
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
