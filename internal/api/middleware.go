package api

import (
	"net/http"
	"strings"

	"github.com/51mans0n/avito-pvz-task/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Считываем заголовок
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"message":"no authorization header"}`, http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"message":"invalid auth header format"}`, http.StatusUnauthorized)
			return
		}
		token := parts[1]
		role, err := auth.ExtractRoleFromToken(token)
		if err != nil {
			http.Error(w, `{"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Кладём role в контекст
		ctx := r.Context()
		ctx = WithRole(ctx, role)

		// передаём дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
