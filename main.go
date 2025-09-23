package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	"github.com/iraqnroll/gochan/controllers"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/templates"
	"github.com/iraqnroll/gochan/views"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	csrfMw := csrf.Protect(
		[]byte("gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"),
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{"localhost:3000"}),
	)

	r.Use(csrfMw)

	// TODO: Make the usage of embedded templates optional.

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	usersC.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "tailwind.gohtml"))
	usersC.Templates.Create = views.Must(views.ParseFS(templates.FS, "createUser.gohtml", "tailwind.gohtml"))

	r.Get("/login", usersC.LoginForm)
	r.Post("/login", usersC.Login)

	r.Post("/logout", usersC.Logout)

	r.Get("/create", usersC.CreateForm)
	r.Post("/create", usersC.Create)

	r.Get("/admin/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", r)
}

//9.9
