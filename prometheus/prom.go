package prometheus

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusClient() (*prometheus.CounterVec, *prometheus.HistogramVec) {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	prometheus.MustRegister(counter)

	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
	prometheus.MustRegister(latencyHistogram)

	return counter, latencyHistogram
}

func PrometheusMiddleware(counter *prometheus.CounterVec, latencyHistogram *prometheus.HistogramVec) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" || c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		//use full path to get /user/:id instead of /user/23
		path := c.FullPath()
		if path == "" {
			//fallback to raw path for 404s or unmatched routes
			path = c.Request.URL.Path
		}
		status := http.StatusText(c.Writer.Status())
		counter.WithLabelValues(c.Request.Method, path, status).Inc()
		// observe latency
		latencyHistogram.WithLabelValues(c.Request.Method, path, status).Observe(float64(time.Since(start).Seconds()))
	}
}

func Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
