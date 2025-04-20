package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "all http"},
		[]string{"method", "path", "status"},
	)

	HttpDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	PVZCreated      = prometheus.NewCounter(prometheus.CounterOpts{Name: "pvz_created_total"})
	ReceptionsAdded = prometheus.NewCounter(prometheus.CounterOpts{Name: "receptions_created_total"})
	ProductsAdded   = prometheus.NewCounter(prometheus.CounterOpts{Name: "products_created_total"})
)

func MustRegister() {
	prometheus.MustRegister(HttpTotal, HttpDur,
		PVZCreated, ReceptionsAdded, ProductsAdded)
}
