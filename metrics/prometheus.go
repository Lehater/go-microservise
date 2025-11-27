// metrics/prometheus.go
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Request duration in seconds",
		},
		[]string{"method", "endpoint"},
	)
)

// InitMetrics регистрирует метрики в глобальном регистре Prometheus.
func InitMetrics() {
	prometheus.MustRegister(totalRequests, requestDuration)
}

// MetricsMiddleware — обёртка над http.Handler, замеряющая RPS и latency.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		totalRequests.WithLabelValues(r.Method, path).Inc()
		next.ServeHTTP(w, r)
		requestDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

// Handler — http.Handler для маршрута /metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}
