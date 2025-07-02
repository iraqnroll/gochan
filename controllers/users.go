package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		Login Template
	}
}

func (u Users) LoginForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Login.Execute(w, nil)
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "<p>Username: %s</p>", r.FormValue("username"))
	fmt.Fprintf(w, "<p>Password: %s</p>", r.FormValue("password"))
}
