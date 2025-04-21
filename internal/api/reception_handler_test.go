package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/model"
)

func (m *mockRepo) CreateReception(ctx context.Context, rec *model.Reception) error {
	args := m.Called(ctx, rec)
	return args.Error(0)
}

func (m *mockRepo) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	args := m.Called(ctx, pvzID)
	rec, _ := args.Get(0).(*model.Reception)
	return rec, args.Error(1)
}

func TestCreateReceptionHandler_Success(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateReceptionHandler(mr)

	mr.On("CreateReception", mock.Anything, mock.AnythingOfType("*model.Reception")).
		Return(nil).
		Once()

	body := `{"pvzId":"31ae2e29-0460-4748-a9f3-2b5747f78960"}`
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString(body))

	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var got model.Reception
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Equal(t, "31ae2e29-0460-4748-a9f3-2b5747f78960", got.PVZID)

	mr.AssertExpectations(t)
}

func TestCloseLastReceptionHandler_Success(t *testing.T) {
	mr := new(mockRepo)
	h := api.CloseLastReceptionHandler(mr)

	r := chi.NewRouter()
	r.Post("/pvz/{pvzId}/close_last_reception", h)

	rec := &model.Reception{
		ID:     "rec-xyz",
		PVZID:  "82cc7cda-bd24-468f-b7b7-844d66b6693c",
		Status: "close",
	}

	mr.On("CloseLastReception", mock.Anything, "82cc7cda-bd24-468f-b7b7-844d66b6693c").
		Return(rec, nil).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/pvz/82cc7cda-bd24-468f-b7b7-844d66b6693c/close_last_reception", nil)
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	mr.AssertExpectations(t)
}

func TestCreateReception_InvalidJSON(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateReceptionHandler(mr)

	req := httptest.NewRequest(http.MethodPost, "/receptions",
		bytes.NewBufferString(`{ invalid json`))
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	mr.AssertExpectations(t)
}

func TestCreateReception_NoPVZ(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateReceptionHandler(mr)

	req := httptest.NewRequest(http.MethodPost, "/receptions",
		bytes.NewBufferString(`{"pvzId":""}`))
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	mr.AssertNotCalled(t, "CreateReception")
}

func TestCreateReception_BadUUID(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateReceptionHandler(mr)

	req := httptest.NewRequest(http.MethodPost, "/receptions",
		bytes.NewBufferString(`{"pvzId":"not‑uuid"}`))
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	mr.AssertNotCalled(t, "CreateReception")
}

func TestCreateReception_RepoErr(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateReceptionHandler(mr)

	mr.On("CreateReception", mock.Anything, mock.Anything).
		Return(assertAnErrorWithMessage("open reception")).Once()

	req := httptest.NewRequest(http.MethodPost, "/receptions",
		bytes.NewBufferString(`{"pvzId":"31ae2e29-0460-4748-a9f3-2b5747f78960"}`))
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	mr.AssertExpectations(t)
}

func TestCloseLastReception_Forbidden(t *testing.T) {
	mr := new(mockRepo)
	h := api.CloseLastReceptionHandler(mr)

	req := httptest.NewRequest(http.MethodPost, "/pvz/111/close_last_reception", nil)
	ctx := api.WithRole(req.Context(), "client")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
	mr.AssertNotCalled(t, "CloseLastReception")
}
