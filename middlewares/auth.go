package middlewares

import (
	"net/http"

	"tincho.dev/rest-ws/server"
)

func AuthMiddleware(s server.Server) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Do something with the request
			next(w, r)
		}
	}
}
