package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)
var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "Duration of gRPC requests in ms",
			Buckets: []float64{50, 100, 250, 1000},
		},
		[]string{"method", "path", "status_code"})
	requestsCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_current_request_count",
			Help: "The current amount of requests which are being handled",
		}, []string{"method", "path"})
)

func init() {
	_ = prometheus.Register(requestDuration)
	_ = prometheus.Register(requestsCurrent)
}

func Prometheus(handler http.Handler, endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL
		method := r.Method
		sw := statusWriter{ResponseWriter: w}

		requestsCurrent.With(prometheus.Labels{
			"method": method,
			"path": path.String(),
		}).Inc()

		defer func(begin time.Time) {
			status := http.StatusText(sw.status)
			duration := time.Since(begin)

			requestDuration.With(prometheus.Labels{
				"method": method,
				"path": endpoint,
				"status_code": status,
			}).Observe(float64(duration.Milliseconds()))

			requestsCurrent.With(prometheus.Labels{
				"method": method,
				"path": endpoint,
			}).Dec()
		}(time.Now())

		handler.ServeHTTP(&sw, r)
	})
}


type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}