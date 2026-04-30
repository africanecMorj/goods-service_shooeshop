package middleware

import (
	"net/http"

	"github.com/africanecMorj/goods-service_shooeshop/internal/domain"
)

func RequirePermission(p domain.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			perms, ok := domain.RolePermissions[role]
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			for _, perm := range perms {
				if perm == p {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

