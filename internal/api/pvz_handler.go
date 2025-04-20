package api

import (
	"encoding/json"
	"fmt"
	"github.com/51mans0n/avito-pvz-task/internal/metrics"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/51mans0n/avito-pvz-task/internal/db"
	"github.com/51mans0n/avito-pvz-task/internal/model"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

// CreatePVZHandler позволяет модератору создавать ПВЗ
func CreatePVZHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "moderator" {
			http.Error(w, `{"message":"access denied"}`, http.StatusForbidden)
			return
		}

		var req struct {
			City string `json:"city"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"message":"invalid json"}`, http.StatusBadRequest)
			return
		}

		if req.City == "" {
			http.Error(w, `{"message":"city is required"}`, http.StatusBadRequest)
			return
		}
		if !allowedCities[req.City] {
			http.Error(w, `{"message":"city not allowed"}`, http.StatusBadRequest)
			return
		}

		pvz := &model.PVZ{
			ID:               uuid.New().String(),
			City:             req.City,
			RegistrationDate: time.Now(),
		}
		if err := repo.CreatePVZ(r.Context(), pvz); err != nil {
			fmt.Printf("Create PVZ error: %v\n", err)
			http.Error(w, `{"message":"server error"}`, http.StatusInternalServerError)
			return
		}

		metrics.PVZCreated.Inc()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pvz)
	}
}

// GetPVZListHandler возвращает список ПВЗ (и их приёмок, товаров) с фильтром и пагинацией
func GetPVZListHandler(repo db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r.Context())
		if role != "moderator" && role != "employee" {
			http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)
			return
		}

		startDateStr := r.URL.Query().Get("startDate")
		endDateStr := r.URL.Query().Get("endDate")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page, limit := parsePageLimit(pageStr, limitStr)

		var startDate, endDate *time.Time
		if startDateStr != "" {
			if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
				startDate = &t
			}
		}
		if endDateStr != "" {
			if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
				endDate = &t
			}
		}

		result, err := repo.GetPVZListWithFilter(r.Context(), startDate, endDate, page, limit)
		if err != nil {
			http.Error(w, `{"message":"server error"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func parsePageLimit(pageStr, limitStr string) (int, int) {
	page := 1
	limit := 10
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 30 {
		limit = l
	}
	return page, limit
}
