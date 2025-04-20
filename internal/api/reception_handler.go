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

// CreateReceptionHandler - employee создаёт приёмку
func CreateReceptionHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "employee" {
			http.Error(w, `{"message":"access denied"}`, http.StatusForbidden)
			return
		}

		var req struct {
			PVZID string `json:"pvzId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if req.PVZID == "" {
			http.Error(w, `{"message":"pvzId is required"}`, http.StatusBadRequest)
			return
		}
		if _, err := uuid.Parse(req.PVZID); err != nil {
			http.Error(w, `{"message":"pvzId invalid format"}`, http.StatusBadRequest)
			return
		}

		rec := &model.Reception{
			ID:       uuid.New().String(),
			PVZID:    req.PVZID,
			DateTime: time.Now(),
			Status:   "in_progress",
		}
		if err := repo.CreateReception(r.Context(), rec); err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		metrics.ReceptionsAdded.Inc()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rec)
	}
}

// CloseLastReceptionHandler - закрытие приёмки
func CloseLastReceptionHandler(repo db.Repository) http.HandlerFunc {
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

		rec, err := repo.CloseLastReception(r.Context(), pvzId)
		if err != nil {
			http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(rec)
	}
}
