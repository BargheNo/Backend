package metricsimpl

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)



type PrometheusMetrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "route", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
	}
}

func (pm *PrometheusMetrics) IncHTTPRequest(method, route, status string) {
	pm.httpRequestsTotal.WithLabelValues(method, route, status).Inc()
}

func (pm *PrometheusMetrics) ObserveHTTPRequestDuration(method, route string, duration float64) {
	pm.httpRequestDuration.WithLabelValues(method, route).Observe(duration)
}
