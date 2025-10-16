package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/controllers"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/templates"
	"github.com/iraqnroll/gochan/views"
)

func main() {
	//1. Setup DB
	cfg := config.InitConfig()
	db, err := config.OpenDBConn(cfg.Database)
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

	//3. Setup reusable page data.
	basePageData := SetupReusablePageData(&boardService, cfg)

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
		BoardService:   &boardService,
		PageData:       basePageData,
	}

	homeC := controllers.Home{
		BoardService:   &boardService,
		GlobalSettings: &cfg.Global,
	}

	boardsC := controllers.Boards{
		BoardService: &boardService,
		PageData:     basePageData,
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
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "head_template.gohtml")), basePageData))

	//Home
	homeC.Home = views.Must(views.ParseFS(templates.FS, "home.gohtml", "head_template.gohtml"))

	//Admin/User
	usersC.Templates.Login = views.Must(views.ParseFS(templates.FS, "login.gohtml", "head_template.gohtml"))
	usersC.Templates.Admin = views.Must(views.ParseFS(templates.FS, "admin.gohtml", "head_template.gohtml"))
	usersC.Templates.Users = views.Must(views.ParseFS(templates.FS, "users.gohtml", "head_template.gohtml"))
	usersC.Templates.Boards = views.Must(views.ParseFS(templates.FS, "boards.gohtml", "head_template.gohtml"))

	//Boards
	boardsC.Board = views.Must(views.ParseFS(templates.FS, "board.gohtml", "head_template.gohtml"))
	boardsC.Thread = views.Must(views.ParseFS(templates.FS, "thread.gohtml", "head_template.gohtml"))

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
		r.Post("/{boardUri}", boardsC.NewThread)

		r.Get("/{boardUri}/{threadId}", boardsC.ThreadForm)
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

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func SetupReusablePageData(bs *models.BoardService, cfg *config.Config) views.BasePageData {
	boards, err := bs.GetBoardList()
	if err != nil {
		panic(err)
	}

	navbar := views.NavbarData{
		BoardList: boards,
	}

	footer := views.FooterData{
		Sitename: cfg.Global.Shortname,
	}

	return views.BasePageData{
		Navbar: navbar,
		Footer: footer,
	}
}
