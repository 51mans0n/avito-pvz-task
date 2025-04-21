package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/stretchr/testify/require"
)

func TestDummyLogin_QueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dummyLogin?role=moderator", nil)
	rr := httptest.NewRecorder()

	api.DummyLoginHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "SOME_TOKEN_moderator", resp["token"])
}

func TestDummyLogin_JSONBody(t *testing.T) {
	body := []byte(`{"role":"employee"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	api.DummyLoginHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "SOME_TOKEN_employee", resp["token"])
}

func TestDummyLogin_NoRole(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dummyLogin", nil)
	rr := httptest.NewRecorder()

	api.DummyLoginHandler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "role is required")
}
