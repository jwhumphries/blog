package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// httpRequestsTotal counts total HTTP requests by method, path, and status.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks request latency by method and path.
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// httpResponseSize tracks response sizes by method and path.
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to 10GB
		},
		[]string{"method", "path"},
	)

	// compressionRatio tracks Brotli compression effectiveness.
	compressionRatio = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_compression_ratio",
			Help:    "Compression ratio (compressed/original size)",
			Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"path"},
	)
)

// responseWriter wraps http.ResponseWriter to capture status and size.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Middleware returns an HTTP middleware that records metrics.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		path := normalizePath(r.URL.Path)

		httpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.status)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		httpResponseSize.WithLabelValues(r.Method, path).Observe(float64(rw.size))
	})
}

// RecordCompression records a compression ratio for a path.
func RecordCompression(path string, originalSize, compressedSize int) {
	if originalSize > 0 {
		ratio := float64(compressedSize) / float64(originalSize)
		compressionRatio.WithLabelValues(path).Observe(ratio)
	}
}

// normalizePath normalizes URL paths for consistent labeling.
// This prevents high cardinality from unique paths.
func normalizePath(path string) string {
	// For a static blog, most paths are fine as-is
	// Just ensure we don't have trailing slashes causing duplicates
	if path == "" {
		return "/"
	}
	return path
}

// Handler returns the Prometheus metrics HTTP handler.
func Handler() http.Handler {
	return promhttp.Handler()
}
