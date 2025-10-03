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
	//1. Setup DB
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

	//2. Setup services
	userService := models.UserService{
		DB: db,
	}

	boardService := models.BoardService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
		BoardService:   &boardService,
	}

	homeC := controllers.Home{
		BoardService: &boardService,
	}

	boardsC := controllers.Boards{
		BoardService: &boardService,
	}

	r := chi.NewRouter()

	//3. Setup Middlewares
	csrfMw := csrf.Protect(
		[]byte("gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"),
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{"localhost:3000"}),
	)

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	// TODO: Make the usage of embedded templates optional.

	//Static
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	//Home
	homeC.Home = views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))

	//Admin/User
	usersC.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "tailwind.gohtml"))
	usersC.Templates.Admin = views.Must(views.ParseFS(templates.FS, "admin.gohtml", "tailwind.gohtml"))
	usersC.Templates.Users = views.Must(views.ParseFS(templates.FS, "users.gohtml", "tailwind.gohtml"))
	usersC.Templates.Boards = views.Must(views.ParseFS(templates.FS, "boards.gohtml", "tailwind.gohtml"))

	//Boards
	boardsC.Board = views.Must(views.ParseFS(templates.FS, "board.gohtml", "tailwind.gohtml"))

	//5. Setup routes

	//Main routing
	r.Route("/", func(r chi.Router) {
		//Home page
		r.Get("/", homeC.HomePage)

		//Login, Logout
		r.Get("/login", usersC.LoginForm)
		r.Post("/login", usersC.Login)
		r.Post("/logout", usersC.Logout)

		//Boards
		r.Get("/{boardUri}", boardsC.BoardForm)
	})

	//Admin panel routes (only accessible to authenticated users)
	r.Route("/admin", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.AdminForm)

		//User management
		r.Get("/users", usersC.UsersForm)
		r.Post("/users/create", usersC.Create)
		r.Post("/users/delete", usersC.Delete)

		//Board management
		r.Get("/boards", usersC.BoardsForm)
		r.Post("/boards/create", usersC.CreateBoard)
		r.Post("/boards/delete", usersC.DeleteBoard)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", r)
}

//9.9
