package api_test

import (
	"bytes"
	"context"
	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (m *mockRepo) CreateProduct(ctx context.Context, pvzID string, prod *model.Product) error {
	args := m.Called(ctx, pvzID, prod)
	return args.Error(0)
}
func (m *mockRepo) DeleteLastProduct(ctx context.Context, pvzID string) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

// helper for error
func assertAnErrorWithMessage(msg string) error {
	return &myFakeError{msg}
}

type myFakeError struct {
	s string
}

func (e *myFakeError) Error() string { return e.s }

func TestCreateProductHandler_Success(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateProductHandler(mr)

	mr.On("CreateProduct", mock.Anything, "82cc7cda-bd24-468f-b7b7-844d66b6693c", mock.AnythingOfType("*model.Product")).
		Return(nil).Once()

	body := `{"type":"электроника","pvzId":"82cc7cda-bd24-468f-b7b7-844d66b6693c"}`
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	mr.AssertExpectations(t)
}

func TestCreateProductHandler_NoActiveReception(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateProductHandler(mr)

	mr.On("CreateProduct", mock.Anything, "82cc7cda-bd24-468f-b7b7-844d66b6693c", mock.AnythingOfType("*model.Product")).
		Return(assertAnErrorWithMessage("no active reception found")).Once()

	body := `{"type":"обувь","pvzId":"82cc7cda-bd24-468f-b7b7-844d66b6693c"}`
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	mr.AssertExpectations(t)
}

func TestCreateProductHandler_Forbidden(t *testing.T) {
	mr := new(mockRepo)
	h := api.CreateProductHandler(mr)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(`{"type":"электроника","pvzId":"82cc7cda-bd24-468f-b7b7-844d66b6693c"}`))
	ctx := api.WithRole(req.Context(), "client")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code)

	mr.AssertNotCalled(t, "CreateProduct")
}

func TestDeleteLastProductHandler_Success(t *testing.T) {
	mr := new(mockRepo)
	h := api.DeleteLastProductHandler(mr)

	mr.On("DeleteLastProduct", mock.Anything, "82cc7cda-bd24-468f-b7b7-844d66b6693c").
		Return(nil).Once()

	r := chi.NewRouter()
	r.Post("/pvz/{pvzId}/delete_last_product", h)

	req := httptest.NewRequest(http.MethodPost, "/pvz/82cc7cda-bd24-468f-b7b7-844d66b6693c/delete_last_product", nil)
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, `{"message":"last product deleted"}`, rr.Body.String())

	mr.AssertExpectations(t)
}

func TestDeleteLastProductHandler_NoProducts(t *testing.T) {
	mr := new(mockRepo)
	h := api.DeleteLastProductHandler(mr)

	mr.On("DeleteLastProduct", mock.Anything, "82cc7cda-bd24-468f-b7b7-844d66b6693c").
		Return(assertAnErrorWithMessage("no products to delete")).Once()

	r := chi.NewRouter()
	r.Post("/pvz/{pvzId}/delete_last_product", h)

	req := httptest.NewRequest(http.MethodPost, "/pvz/82cc7cda-bd24-468f-b7b7-844d66b6693c/delete_last_product", nil)
	ctx := api.WithRole(req.Context(), "employee")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "no products to delete")

	mr.AssertExpectations(t)
}
