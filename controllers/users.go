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
	u.Templates.Login.Execute(w, nil)
}

func (u Users) CreateForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Create.Execute(w, nil)
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "<p>Username: %s</p>", r.FormValue("username"))
	fmt.Fprintf(w, "<p>Password: %s</p>", r.FormValue("password"))
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
