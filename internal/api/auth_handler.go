package api

import (
	"encoding/json"
	"net/http"
)

type DummyLoginRequest struct {
	Role string `json:"role"`
}

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"invalid json"}`))
		return
	}

	token := "SOME_TOKEN_" + req.Role
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"` + token + `"`))
}
