package api

import (
	"encoding/json"
	"github.com/51mans0n/avito-pvz-task/internal/db"
	"github.com/51mans0n/avito-pvz-task/internal/metrics"
	"github.com/51mans0n/avito-pvz-task/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// CreateProductHandler - employee добавляет товар
func CreateProductHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "employee" {
			http.Error(w, `{"message":"access denied"}`, http.StatusForbidden)
			return
		}

		var req struct {
			Type  string `json:"type"`
			PVZID string `json:"pvzId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if req.Type == "" || req.PVZID == "" {
			http.Error(w, `{"message":"type and pvzId are required"}`, http.StatusBadRequest)
			return
		}
		if _, err := uuid.Parse(req.PVZID); err != nil {
			http.Error(w, `{"message":"pvzId invalid"}`, http.StatusBadRequest)
			return
		}

		prod := &model.Product{
			ID:       uuid.New().String(),
			Type:     req.Type,
			DateTime: time.Now(),
		}

		if err := repo.CreateProduct(r.Context(), req.PVZID, prod); err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		metrics.ProductsAdded.Inc()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(prod)
	}
}

// DeleteLastProductHandler - удалить последний товар LIFO
func DeleteLastProductHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "employee" {
			http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)
			return
		}

		pvzId := chi.URLParam(r, "pvzId")
		if pvzId == "" {
			http.Error(w, `{"message":"invalid pvzId"}`, http.StatusBadRequest)
			return
		}

		if err := repo.DeleteLastProduct(r.Context(), pvzId); err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"last product deleted"}`))
	}
}
