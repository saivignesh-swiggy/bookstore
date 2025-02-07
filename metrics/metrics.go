package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// APIMetrics struct holds Prometheus metrics
type APIMetrics struct {
	TotalBooksAdded *prometheus.CounterVec
	BooksAvailable  prometheus.Gauge
	HTTPRequests    *prometheus.CounterVec
	HTTPDuration    *prometheus.HistogramVec
	QueryCount      *prometheus.CounterVec
	QueryDuration   *prometheus.HistogramVec
}

var (
	once       sync.Once
	apiMetrics *APIMetrics
)

// NewAPIMetrics initializes Prometheus metrics
func NewAPIMetrics() *APIMetrics {
	once.Do(func() {
		apiMetrics = &APIMetrics{
			TotalBooksAdded: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "total_books_added",
					Help: "Total number of books added to the bookstore.",
				},
				[]string{"status"},
			),
			BooksAvailable: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: "books_available",
					Help: "The number of books available in the bookstore.",
				},
			),
			HTTPRequests: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Total number of HTTP requests made.",
				},
				[]string{"action", "status"},
			),
			HTTPDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_duration_seconds",
					Help:    "Histogram of HTTP request durations.",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"action"},
			),
			QueryCount: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "query_count_total",
					Help: "Total number of database queries executed.",
				},
				[]string{"query_type"},
			),
			QueryDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "query_duration_seconds",
					Help:    "Histogram of database query durations.",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"query_type"},
			),
		}

		// Register all metrics only once
		prometheus.MustRegister(
			apiMetrics.TotalBooksAdded,
			apiMetrics.BooksAvailable,
			apiMetrics.HTTPRequests,
			apiMetrics.HTTPDuration,
			apiMetrics.QueryCount,
			apiMetrics.QueryDuration,
		)
	})

	return apiMetrics
}

// HTTP Metrics
func (m *APIMetrics) IncRequestCount(method, status string) {
	m.HTTPRequests.WithLabelValues(method, status).Inc()
}

func (m *APIMetrics) MeasureRequestLatency(method string, duration float64) {
	m.HTTPDuration.WithLabelValues(method).Observe(duration)
}

// Query Metrics
func (m *APIMetrics) IncQueryCount(queryType string) {
	m.QueryCount.WithLabelValues(queryType).Inc()
}

func (m *APIMetrics) MeasureQueryLatency(queryType string, duration float64) {
	m.QueryDuration.WithLabelValues(queryType).Observe(duration)
}

// Book Metrics
func (m *APIMetrics) IncBooksAdded(status string) {
	m.TotalBooksAdded.WithLabelValues(status).Inc()
}

func (m *APIMetrics) UpdateBooksAvailable(count int) {
	m.BooksAvailable.Set(float64(count))
}

// Expose Metrics Endpoint
func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9090", nil)
}
