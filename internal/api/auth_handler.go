package api

import (
	"encoding/json"
	"github.com/51mans0n/avito-pvz-task/internal/auth"
	"net/http"
)

// DummyLoginHandler ­— единственная ручка, которая выдаёт тестовый токен.
// Теперь роль читаем НЕ из JSON‑тела, а из query‑параметра ?role=…
func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	if role == "" {
		http.Error(w, `{"message":"role query param required"}`, http.StatusBadRequest)
		return
	}

	token := auth.IssueDummyToken(role) // вернёт строку вида SOME_TOKEN_<role>
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
