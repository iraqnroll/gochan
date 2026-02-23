package middlewares

import (
	"fmt"
	"net/http"

	"github.com/iraqnroll/gochan/cmd/front/handlers"
	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/db/services"
)

type UsersMiddleware struct {
	SessionService *services.SessionService
}

func NewUsersMiddleware(sessionSvc *services.SessionService) (umw UsersMiddleware) {
	umw.SessionService = sessionSvc
	return umw
}

func (umw UsersMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := handlers.ReadCookie(r, handlers.CookieSession)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.GetUser(tokenCookie)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		//If we get to this point, we have an active user - store it in the context.
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
