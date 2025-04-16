package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/db"
	"github.com/51mans0n/avito-pvz-task/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockRepo реализует db.Repository
type mockRepo struct {
	mock.Mock
}

var _ db.Repository = (*mockRepo)(nil)

func (m *mockRepo) CreatePVZ(ctx context.Context, pvz *model.PVZ) error {
	args := m.Called(ctx, pvz)
	return args.Error(0)
}
func (m *mockRepo) GetPVZListWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]model.PVZWithReceptions, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]model.PVZWithReceptions), args.Error(1)
}

func TestCreatePVZHandler_Success(t *testing.T) {
	mrepo := new(mockRepo)
	h := api.CreatePVZHandler(mrepo)

	mrepo.On("CreatePVZ", mock.Anything, mock.AnythingOfType("*model.PVZ")).
		Return(nil).
		Once()

	body := `{"city":"Москва"}`
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := api.WithRole(req.Context(), "moderator")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	var pvz model.PVZ
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &pvz))
	require.Equal(t, "Москва", pvz.City)
	require.NotEmpty(t, pvz.ID)

	mrepo.AssertExpectations(t)
}

func TestCreatePVZHandler_Forbidden(t *testing.T) {
	mrepo := new(mockRepo)
	h := api.CreatePVZHandler(mrepo)

	body := `{"city":"Казань"}`
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))

	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
	mrepo.AssertNotCalled(t, "CreatePVZ")
}

func TestCreatePVZHandler_CityNotAllowed(t *testing.T) {
	mrepo := new(mockRepo)
	h := api.CreatePVZHandler(mrepo)

	body := `{"city":"Новосибирск"}`
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
	ctx := api.WithRole(req.Context(), "moderator")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	mrepo.AssertNotCalled(t, "CreatePVZ")
}

func TestCreatePVZHandler_RepoError(t *testing.T) {
	mrepo := new(mockRepo)
	h := api.CreatePVZHandler(mrepo)

	mrepo.On("CreatePVZ", mock.Anything, mock.AnythingOfType("*model.PVZ")).
		Return(assertAnErrorWithMessage("some db error")).
		Once()

	body := `{"city":"Москва"}`
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := api.WithRole(req.Context(), "moderator")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	mrepo.AssertExpectations(t)
}

func TestGetPVZListHandler_Success(t *testing.T) {
	mr := new(mockRepo)
	h := api.GetPVZListHandler(mr)

	result := []model.PVZWithReceptions{
		{
			PVZ: &model.PVZResponse{
				ID:   "82cc7cda-bd24-468f-b7b7-844d66b6693c",
				City: "Москва",
			},
			Receptions: []model.ReceptionWithProd{
				{
					Reception: &model.ReceptionResponse{ID: "rec-123", Status: "in_progress"},
					Products:  []model.ProductResponse{{ID: "prod-1", Type: "электроника"}},
				},
			},
		},
	}

	mr.On("GetPVZListWithFilter", mock.Anything, (*time.Time)(nil), (*time.Time)(nil), 1, 10).
		Return(result, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/pvz?page=1&limit=10", nil)
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var arr []model.PVZWithReceptions
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &arr))
	require.Len(t, arr, 1)
	require.Equal(t, "82cc7cda-bd24-468f-b7b7-844d66b6693c", arr[0].PVZ.ID)

	mr.AssertExpectations(t)
}

func TestGetPVZListHandler_Forbidden(t *testing.T) {
	mr := new(mockRepo)
	h := api.GetPVZListHandler(mr)

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	ctx := api.WithRole(req.Context(), "client")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
	mr.AssertNotCalled(t, "GetPVZListWithFilter")
}

func TestGetPVZListHandler_RepoError(t *testing.T) {
	mr := new(mockRepo)
	h := api.GetPVZListHandler(mr)

	mr.On("GetPVZListWithFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]model.PVZWithReceptions(nil), assertAnErrorWithMessage("db error")).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	ctx := api.WithRole(req.Context(), "moderator")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	mr.AssertExpectations(t)
}
