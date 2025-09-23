package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/iraqnroll/gochan/models"
)

type Users struct {
	Templates struct {
		Login  Template
		Create Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) LoginForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Login.Execute(w, r, nil)
}

func (u Users) CreateForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Create.Execute(w, r, nil)
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
	}

	var data struct {
		username string
		password string
	}
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")

	//Try to authenticate
	user, err := u.UserService.Authenticate(data.username, data.password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	//Create a new active session on login
	session, err := u.SessionService.CreateSession(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went horribly wrong...", http.StatusInternalServerError)
		return
	}

	//Authentication successful, add the Set-Cookie header before writting the response
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/admin/me", http.StatusFound)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	user_type, usert_err := strconv.Atoi(r.FormValue("user_type"))
	if usert_err != nil {
		http.Error(w, "Invalid user_type entry, only numeric values are allowed.", http.StatusBadRequest)
		return
	}

	user, err := u.UserService.Create(username, password, email, user_type)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went horribly wrong....", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	user, err := u.SessionService.GetUser(tokenCookie)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %s", user.Username)
}

func (u Users) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = u.SessionService.DeleteSession(token)
	if err != nil {
		fmt.Println("Logout failed : %w", err)
		http.Error(w, "Something went horribly wrong...", http.StatusInternalServerError)
		return
	}

	//Delete the session cookie and redirect to homepage
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/", http.StatusFound)
}
