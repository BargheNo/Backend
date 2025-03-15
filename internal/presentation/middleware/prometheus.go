package middleware

import (
	"strconv"
	"time"

	"github.com/BargheNo/Backend/internal/domain/metrics"
	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMiddleware struct{
	metricsClient metrics.MetricsClient
}

func NewPrometheusMiddleware(metricsClient metrics.MetricsClient) *PrometheusMiddleware {
	return &PrometheusMiddleware{
		metricsClient: metricsClient,
	}
}

func (pm *PrometheusMiddleware) PrometheusMiddleware(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds()
	status := c.Writer.Status()

	route := c.Request.URL.Path
	method := c.Request.Method

	pm.metricsClient.IncHTTPRequest(method, route, strconv.Itoa(status))
    pm.metricsClient.ObserveHTTPRequestDuration(method, route, duration)
}

func SetupPrometheusEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
