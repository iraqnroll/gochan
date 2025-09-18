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
	UserService *models.UserService
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

	//Authentication successful, add the Set-Cookie header before writting the response
	cookie := http.Cookie{
		Name:     "username",
		Value:    user.Username,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User authenticated: %v", user)
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
	username, err := r.Cookie("username")
	if err != nil {
		fmt.Fprintf(w, "The auth cookie could not be read.")
		return
	}

	fmt.Fprintf(w, "Auth cookie: %s", username.Value)
}
