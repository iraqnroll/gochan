package middlewares

import (
	"fmt"
	"net/http"

	"github.com/iraqnroll/gochan/cmd/front/handlers"
	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/rand"
)

func SetFingerprint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if user already has a generated fingerprint in cookies
		_, err := handlers.ReadCookie(r, handlers.UserFingerprint)

		if err != nil {
			userIp := handlers.GetClientIp(r)
			fingerprint := rand.GenerateFingerprint(userIp, "ZX7Q9M2K")
			fmt.Printf("[Fingerprint] Generated new fingerprint for %s - %s", userIp, fingerprint)

			ctx := r.Context()
			ctx = context.StoreFingerprint(ctx, fingerprint)
			r = r.WithContext(ctx)

			handlers.SetCookie(w, handlers.UserFingerprint, fingerprint)
			next.ServeHTTP(w, r)
		} else {
			//Fingerprint already present, skip generating.
			next.ServeHTTP(w, r)
			return
		}
	})
}

func RequireFingerprint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fingerprint := context.GetFingerprint(r.Context())
		if fingerprint == "" {
			http.Error(w, "This action requires a valid fingerprint.", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
