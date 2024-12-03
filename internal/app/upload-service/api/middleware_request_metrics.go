package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricRequestsTotal       = "http_requests_total"
	metricRequestsTotalHelp   = "Total number of HTTP requests"
	metricRequestDuration     = "http_request_duration_seconds"
	metricRequestDurationHelp = "Histogram of HTTP request durations"
)

func MiddlewareRequestMetrics() gin.HandlerFunc {
	// Define the counter metric for the total number of requests
	metricRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricRequestsTotal,
			Help: metricRequestsTotalHelp,
		},
		[]string{"path", "method"},
	)

	// Define the histogram metric for the duration of requests
	metricRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    metricRequestDuration,
			Help:    metricRequestDurationHelp,
			Buckets: prometheus.DefBuckets, // Default buckets for request duration
		},
		[]string{"path", "method"},
	)

	// Register the metrics with Prometheus
	prometheus.MustRegister(metricRequestsTotal)
	prometheus.MustRegister(metricRequestDuration)

	// Return the middleware function
	return func(c *gin.Context) {
		// Record the start time for duration measurement
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate the request duration
		duration := time.Since(start).Seconds()

		// Record the total number of requests
		metricRequestsTotal.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()

		// Record the request duration in the histogram
		metricRequestDuration.WithLabelValues(c.Request.URL.Path, c.Request.Method).Observe(duration)
	}
}
