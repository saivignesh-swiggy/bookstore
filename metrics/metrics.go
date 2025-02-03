package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	TotalBooksAdded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_books_added",
			Help: "Total number of books added to the bookstore.",
		},
		[]string{"status"},
	)
	BooksAvailable = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "books_available",
			Help: "The number of books available in the bookstore.",
		},
	)
)

func Init() {
	prometheus.MustRegister(TotalBooksAdded)
	prometheus.MustRegister(BooksAvailable)
}

func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
}
