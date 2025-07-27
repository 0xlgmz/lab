package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "path"},
	)

	// New metrics
	httpRequestsByMethod = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_by_method_total",
			Help: "Total number of HTTP requests by method",
		},
		[]string{"service", "method"},
	)

	httpRequestsByStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_by_status_total",
			Help: "Total number of HTTP requests by status code category",
		},
		[]string{"service", "status_category"},
	)

	httpRequestErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"service", "error_type"},
	)
)

// MetricsMiddleware records HTTP metrics for Prometheus
func MetricsMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "undefined"
		}

		c.Next()

		status := c.Writer.Status()
		statusStr := strconv.Itoa(status)
		duration := time.Since(start).Seconds()

		// Record basic metrics
		httpRequestsTotal.WithLabelValues(service, c.Request.Method, path, statusStr).Inc()
		httpRequestDuration.WithLabelValues(service, c.Request.Method, path).Observe(duration)

		// Record method-based metrics
		httpRequestsByMethod.WithLabelValues(service, c.Request.Method).Inc()

		// Record status-based metrics
		statusCategory := getStatusCategory(status)
		httpRequestsByStatus.WithLabelValues(service, statusCategory).Inc()

		// Record error metrics if applicable
		if status >= 400 {
			errorType := getErrorType(status)
			httpRequestErrors.WithLabelValues(service, errorType).Inc()
		}
	}
}

// getStatusCategory returns the category of the HTTP status code
func getStatusCategory(status int) string {
	switch {
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	case status >= 300:
		return "redirect"
	case status >= 200:
		return "success"
	case status >= 100:
		return "informational"
	default:
		return "unknown"
	}
}

// getErrorType returns the type of error based on the status code
func getErrorType(status int) string {
	switch {
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	default:
		return "unknown"
	}
}
