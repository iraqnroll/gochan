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
	"github.com/gorilla/csrf"
	"github.com/iraqnroll/gochan/cmd/front/handlers"
	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/repos"
	"github.com/iraqnroll/gochan/services"
	"github.com/iraqnroll/gochan/views"
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

	a.Router = chi.NewRouter()
	a.InitMiddlewares(
		cfg.Frontend.CsrfKey,
		cfg.Frontend.CsrfSecure,
		cfg.Frontend.CsrfTrustedOrigins,
	)

	a.InitServices()
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
	a.FileService = services.NewFileService()
	a.UserService = services.NewUserService(uRepo)
	a.SessionService = services.NewSessionService(sRepo, a.Settings.Api.SessionTokenSize)

	a.ThreadService = services.NewThreadService(tRepo, a.PostService)
	a.BoardService = services.NewBoardService(bRepo, a.ThreadService, a.FileService)
}

func (a *Frontend) InitRoutes() {
	footerData := &views.FooterData{Sitename: a.Settings.Global.Shortname}
	homeHandler := handlers.Home{}
	homeHandler.Footer = footerData

	//Base route
	a.Router.Route("/", func(r chi.Router) {
		r.Get("/", homeHandler.Home)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Page not found.", http.StatusNotFound)
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

// TODO: add usermiddleware once i refactor route handlers for frontend
func (a *Frontend) InitMiddlewares(csrfKey string, csrfSecure bool, trusted_origins []string) {
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(csrfSecure),
		csrf.TrustedOrigins(trusted_origins),
	)

	a.Router.Use(csrfMw)
	a.Router.Use(middleware.Logger)
}
