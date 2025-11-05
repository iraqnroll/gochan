package board

import "net/http"

type API struct{}

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
func (a *API) List(w http.ResponseWriter, r *http.Request)   {}
func (a *API) Create(w http.ResponseWriter, r *http.Request) {}
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {}
func (a *API) Update(w http.ResponseWriter, r *http.Request) {}
func (a *API) Get(w http.ResponseWriter, r *http.Request)    {}
