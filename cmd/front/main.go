package front

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iraqnroll/gochan/cmd/front/handlers"
	"github.com/iraqnroll/gochan/cmd/front/middlewares"
	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/repos"
	"github.com/iraqnroll/gochan/services"
	"github.com/pressly/goose"
)

type Frontend struct {
	DB       *sql.DB
	Router   *chi.Mux
	Settings *config.Config

	BoardService   *services.BoardService
	PostService    *services.PostService
	ThreadService  *services.ThreadService
	FileService    *services.FileService
	SessionService *services.SessionService
	UserService    *services.UserService
}

func (a *Frontend) Run(host string) {
	fmt.Println("Frontend running on host - ", host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *Frontend) Init(cfg *config.Config) {
	//TODO: Refactor this shit, config that takes config as a parameter ?
	db, err := config.OpenDBConn(cfg.Database)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	a.Settings = cfg
	a.DB = db

	goose.SetLogger(log.Default())
	if err := goose.Up(db, "./migrations"); err != nil {
		panic(err)
	}
	fmt.Println("Migrations applied successfully.")

	a.Router = chi.NewRouter()
	a.InitServices()
	a.InitMiddlewares()
	a.InitFileServer(cfg.Frontend.StaticDir, "static")
	a.InitRoutes()
}

//TODO: Repository object lifetime
/*
	I dont know if this is a good approach...
	Repository objects will live only while services depend on them, if we're
	thinking about this from the microservice perspective, if one service fails - the whole pipeline
	down to the data access layer gets destroyed, we'd have to re-initialize the whole app (services and up).

	If I end up with a substantial amount of independant services I might refactor this.
*/
func (a *Frontend) InitServices() {
	bRepo := repos.NewPostgresBoardRepository(a.DB)
	pRepo := repos.NewPostgresPostRepository(a.DB)
	tRepo := repos.NewPostgresThreadRepository(a.DB)
	uRepo := repos.NewPostgresUserRepository(a.DB)
	sRepo := repos.NewPostgresSessionRepository(a.DB)

	a.PostService = services.NewPostService(pRepo)
	a.FileService = services.NewFileService(a.Settings.Global.AllowedMediaTypes)
	a.UserService = services.NewUserService(uRepo)
	a.SessionService = services.NewSessionService(sRepo, a.Settings.Api.SessionTokenSize)

	a.ThreadService = services.NewThreadService(tRepo, a.PostService)
	a.BoardService = services.NewBoardService(bRepo, a.ThreadService, a.FileService)
}

// TODO: Implement caching of viewmodels/global page data
func (a *Frontend) InitRoutes() {
	footerData := models.FooterData{Sitename: a.Settings.Global.Shortname}
	parentPageData := models.ParentPageData{
		Footer:    footerData,
		Shortname: a.Settings.Global.Shortname,
		Subtitle:  a.Settings.Global.Subtitle}

	homeHandler := handlers.NewHomeHandler(a.BoardService, a.PostService, parentPageData, a.Settings.Global.RecentPostsNum)
	boardHandler := handlers.NewBoardsHandler(a.BoardService, a.ThreadService, a.FileService, parentPageData, 10)
	threadHandler := handlers.NewThreadsHandler(a.ThreadService, a.PostService, a.FileService, parentPageData, 50)
	usersHandler := handlers.NewUsersHandler(a.UserService, a.SessionService, parentPageData)
	modHandler := handlers.NewModHandler(a.UserService, a.BoardService, a.ThreadService, a.FileService, parentPageData)

	//Base route
	a.Router.Route("/", func(r chi.Router) {
		r.Get("/", homeHandler.Home)
		r.Get("/login", usersHandler.LoginPage)
		r.Post("/login", usersHandler.Login)
		r.Post("/logout", usersHandler.Logout)

		r.Get("/{board_uri}", boardHandler.Board)
		r.Post("/{board_uri}", boardHandler.NewThread)

		r.Get("/{board_uri}/{thread_id}", threadHandler.Thread)
		r.Post("/{board_uri}/{thread_id}", threadHandler.Reply)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Page not found.", http.StatusNotFound)
		})
	})

	a.Router.Route("/mod", func(r chi.Router) {
		r.Use(middlewares.RequireUser)
		r.Get("/", modHandler.ModPage)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", modHandler.ModUsers)
			r.Get("/update/{user_id}", modHandler.EditUserPage)

			r.Post("/create", modHandler.CreateUser)
			r.Post("/update/{user_id}", modHandler.UpdateUser)
			r.Post("/delete/{user_id}", modHandler.DeleteUser)
		})
		r.Route("/boards", func(r chi.Router) {
			r.Get("/", modHandler.ModBoards)
			r.Get("/update/{board_id}", modHandler.EditBoardPage)

			r.Post("/create", modHandler.CreateBoard)
			r.Post("/delete", modHandler.DeleteBoard)
			r.Post("/update/{board_id}", modHandler.UpdateBoard)

		})
	})
}

func (a *Frontend) InitFileServer(path string, workdir string) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, workdir))

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		a.Router.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	a.Router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})
}

// TODO: gorilla.csrf is deprecated, i removed csrf token from all forms, will need to re-implement later...
func (a *Frontend) InitMiddlewares() {
	userMw := middlewares.NewUsersMiddleware(a.SessionService)

	a.Router.Use(middleware.Logger)
	a.Router.Use(userMw.SetUser)
}
