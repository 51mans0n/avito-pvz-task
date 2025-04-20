package metrics

import (
	"net/http"
	"strconv"
	"time"
)

func PromMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wr{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		path := r.URL.Path // для MVP без шаблонов
		HttpTotal.WithLabelValues(r.Method, path, strconv.Itoa(ww.status)).Inc()
		HttpDur.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

type wr struct {
	http.ResponseWriter
	status int
}

func (w *wr) WriteHeader(c int) { w.status = c; w.ResponseWriter.WriteHeader(c) }
