package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const RoleKey contextKey = "role"

func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				http.Error(w, "unauthorized", 401)
				return
			}

			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", 401)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "unauthorized", 401)
				return
			}

			role, _ := claims["role"].(string)

			ctx := context.WithValue(r.Context(), RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}