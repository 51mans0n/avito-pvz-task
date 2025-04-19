package api

import (
	"net/http"
	"strings"

	"github.com/51mans0n/avito-pvz-task/internal/auth"
)

// AuthMiddleware проверяет Bearer‑токен и вкладывает роль в контекст.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, `missing bearer token`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		role, err := auth.ExtractRole(token) // <-- auth.ExtractRole мы писали раньше
		if err != nil {
			http.Error(w, `unauthorized: `+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := WithRole(r.Context(), role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
