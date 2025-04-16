package api

import (
	"encoding/json"
	"net/http"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Role string `json:"role"`
	}
	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		http.Error(w, `{"message":"role is required"}`, http.StatusBadRequest)
		return
	}

	token := "SOME_TOKEN_" + req.Role
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + token + `"}`))
}
