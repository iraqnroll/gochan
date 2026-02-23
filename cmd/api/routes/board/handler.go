package board

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/db/services"
)

type API struct {
	BoardService *services.BoardService
}

// List godoc
//
//	@summary        List boards
//	@description    List all registered boards
//	@tags           boards
//	@accept         json
//	@produce        json
//	@success        200 {array}     []models.BoardDto
//	@failure        500 {object}    err.Error
//	@router         /boards [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	result, err := a.BoardService.GetBoardList()
	if err != nil {
		fmt.Printf("Error : %s", err)
		http.Error(w, "Error while retrieving the boards", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error while converting response to JSON", http.StatusInternalServerError)
		return
	}
}

// Get godoc
//
//	@summary        Gets a specific board by URI
//	@description    Gets a specific active boards by its URI
//	@tags           boards
//	@accept         json
//	@produce        json
//	@success        200 {object}     models.BoardDto
//	@failure        500 {object}    err.Error
//	@router         /boards/{uri} [get]
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	//TODO: add parameter validation
	uri := chi.URLParam(r, "uri")

	result, err := a.BoardService.GetBoard(uri, false)
	if err != nil {
		http.Error(w, "Failed to fetch board by uri", http.StatusInternalServerError)
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error while converting response to JSON", http.StatusInternalServerError)
		return
	}
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {}
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {}
func (a *API) Update(w http.ResponseWriter, r *http.Request) {}
