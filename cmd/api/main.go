package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iraqnroll/gochan/cmd/api/routes/board"
	"github.com/iraqnroll/gochan/cmd/api/routes/post"
	"github.com/iraqnroll/gochan/cmd/api/routes/thread"
	"github.com/iraqnroll/gochan/cmd/api/routes/user"
	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/db/repos"
	"github.com/iraqnroll/gochan/db/services"
)

type Api struct {
	DB     *sql.DB
	Router *chi.Mux

	BoardService   *services.BoardService
	PostService    *services.PostService
	ThreadService  *services.ThreadService
	FileService    *services.FileService
	SessionService *services.SessionService
	UserService    *services.UserService
}

// @title          gochan API
// @version        1.0
// @description    RESTful API for interaction with gochan backend

// @contact.name   Lukas T.
// @contact.url    https://likeadaydream.lt

// @host       localhost:3000
// @basePath   /api
func (a *Api) Run(host string) {
	fmt.Println("API running on host - ", host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *Api) Init() {
	//TODO: Refactor this shit, config that takes config as a parameter ?
	db, err := config.OpenDBConn()
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	a.DB = db
	a.Router = chi.NewRouter()

	a.InitServices()
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

func (a *Api) InitServices() {
	bRepo := repos.NewPostgresBoardRepository(a.DB)
	pRepo := repos.NewPostgresPostRepository(a.DB)
	tRepo := repos.NewPostgresThreadRepository(a.DB)
	uRepo := repos.NewPostgresUserRepository(a.DB)
	sRepo := repos.NewPostgresSessionRepository(a.DB)

	a.PostService = services.NewPostService(pRepo, config.FingerprintSalt(), config.TripcodeSalt())
	a.FileService = services.NewFileService(config.AllowedMediaTypes())
	a.UserService = services.NewUserService(uRepo)
	a.SessionService = services.NewSessionService(sRepo, config.SessionTokenSize())

	a.ThreadService = services.NewThreadService(tRepo, a.PostService)
	a.BoardService = services.NewBoardService(bRepo, a.ThreadService, a.FileService)
}

func (a *Api) InitRoutes() {
	postAPI := &post.API{PostService: a.PostService, RecentPostsNum: config.NumberOfRecentPosts()}
	threadAPI := &thread.API{ThreadService: a.ThreadService}
	boardAPI := &board.API{BoardService: a.BoardService}
	userAPI := &user.API{}

	//Base route
	a.Router.Route("/api", func(r chi.Router) {
		r.Get("/boards", boardAPI.List)

		r.Get("/threads/{id}", threadAPI.Get)
		r.Get("/boards/{uri}", boardAPI.Get)

		r.Post("/threads", threadAPI.Create)
		r.Post("/login", userAPI.Login)

		//Route for outhorized users.
		r.Route("/user", func(r chi.Router) {
			r.Get("/all", userAPI.List)
			r.Get("/{id}", userAPI.Get)
			r.Post("/create", userAPI.Create)
			r.Post("/update/{id}", userAPI.Update)

			r.Route("/boards", func(r chi.Router) {
				r.Post("/create", boardAPI.Create)
				r.Post("/delete/{id}", boardAPI.Delete)
				r.Post("/update/{id}", boardAPI.Update)
			})
		})

		r.Route("/posts", func(r chi.Router) {
			r.Get("/mostrecent", postAPI.ListMostRecent)
			r.Post("/", postAPI.Create)
			r.Post("/{id}", postAPI.Update)
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Page not found.", http.StatusNotFound)
		})
	})
}
