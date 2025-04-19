package api

import (
	"encoding/json"
	"github.com/51mans0n/avito-pvz-task/internal/auth"
	"net/http"
)

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")

	// если query‑парам нет – пробуем JSON‑тело {"role":"…"}
	if role == "" && r.Method == http.MethodPost {
		var body struct {
			Role string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			role = body.Role
		}
	}

	if role == "" {
		http.Error(w, `{"message":"role is required"}`, http.StatusBadRequest)
		return
	}

	token := auth.IssueDummyToken(role)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
