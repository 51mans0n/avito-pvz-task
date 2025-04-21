package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/51mans0n/avito-pvz-task/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestPromMiddleware_Increments(t *testing.T) {
	metrics.MustRegister()

	h := metrics.PromMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusTeapot, rr.Code)

	total := testutil.ToFloat64(
		metrics.HttpTotal.WithLabelValues("GET", "/test", "418"),
	)
	require.Equal(t, 1.0, total)
}
