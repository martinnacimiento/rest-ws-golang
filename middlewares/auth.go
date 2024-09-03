package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"tincho.dev/rest-ws/models"
	"tincho.dev/rest-ws/server"
)

type contextKey string

const ClaimsKey contextKey = "claims"

var (
	NO_AUTH_ROUTES = []string{"/", "/signup", "/signin", "/users"}
)

func shouldAuth(route string) bool {
	for _, r := range NO_AUTH_ROUTES {
		if strings.Contains(r, route) {
			return false
		}
	}

	return true
}

func AuthMiddleware(s server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldAuth(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authorizationToken := r.Header.Get("Authorization")
			tokenString := strings.Replace(authorizationToken, "Bearer ", "", 1)

			token, err := jwt.ParseWithClaims(
				tokenString,
				&models.AppClaims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(s.Config().JWTSecret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*models.AppClaims)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
