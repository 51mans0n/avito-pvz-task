package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_OK(t *testing.T) {
	// handler‑эхо проверит, что роль попала в context
	h := api.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := api.GetRole(r.Context())
		if role == "employee" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer SOME_TOKEN_employee")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	h := api.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
