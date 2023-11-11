package web

import (
	"net/http"

	"documents/internal/core"
)

func CORSMiddleware(repo *core.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", repo.Config.Server.CORSOrigin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin,access-control-allow-headers")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
