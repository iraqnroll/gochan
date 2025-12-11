package handlers

import (
	"fmt"
	"net/http"

	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
)

const (
	CookieSession = "gochan.session"
)

type Users struct {
	UserService    *services.UserService
	SessionService *services.SessionService
	ParentPage     models.ParentPageData
}

func NewUsersHandler(userSvc *services.UserService, sessionSvc *services.SessionService, parentPage models.ParentPageData) (u Users) {
	u.UserService = userSvc
	u.SessionService = sessionSvc
	u.ParentPage = parentPage

	return u
}

func (u Users) LoginPage(w http.ResponseWriter, r *http.Request) {
	views.LoginPageComponent(u.ParentPage).Render(r.Context(), w)
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var data struct {
		Username string
		Password string
	}

	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")
	//TODO: Add login data validation
	user, err := u.UserService.Authenticate(data.Username, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//User authenticated, create a new active session
	session, err := u.SessionService.CreateNew(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("generated sesh token: %s\n", session.Token)
	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (u Users) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := ReadCookie(r, CookieSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = u.SessionService.DeleteSession(token)
	if err != nil {
		http.Error(w, "Something went horribly wrong..."+err.Error(), http.StatusInternalServerError)
		return
	}
	DeleteCookie(w, CookieSession)
	http.Redirect(w, r, "/", http.StatusFound)
}

func NewCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookie
}

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := NewCookie(name, value)
	http.SetCookie(w, cookie)
}

func ReadCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := NewCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
