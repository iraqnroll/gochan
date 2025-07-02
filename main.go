package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iraqnroll/gochan/controllers"
	"github.com/iraqnroll/gochan/templates"
	"github.com/iraqnroll/gochan/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// TODO: Make the usage of embedded templates optional.

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	var usersC controllers.Users
	usersC.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "tailwind.gohtml"))

	r.Get("/login", usersC.LoginForm)
	r.Post("/login", usersC.Login)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

//7.8
