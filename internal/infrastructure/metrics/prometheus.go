package metrics

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	metrics             *bootstrap.Metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewPrometheusMetrics(metrics *bootstrap.Metrics) *PrometheusMetrics {
	return &PrometheusMetrics{
		metrics: metrics,
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{

				Name: metrics.HTTPRequestsTotal.Name,
				Help: metrics.HTTPRequestsTotal.Help,
			},
			[]string{"method", "route", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: metrics.HTTPRequestDuration.Name,
				Help: metrics.HTTPRequestDuration.Help,
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
