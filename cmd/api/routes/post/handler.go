package post

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/services"
)

type API struct {
	PostService *services.PostService
	Settings    *config.API
}

func (a *API) ListMostRecent(w http.ResponseWriter, r *http.Request) {
	result, err := a.PostService.GetMostRecent(a.Settings.RecentPostsNum)
	if err != nil {
		fmt.Println("Error: %s \n", err)
		http.Error(w, "Failed to fetch most recent posts", http.StatusInternalServerError)
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
func (a *API) Get(w http.ResponseWriter, r *http.Request)    {}
