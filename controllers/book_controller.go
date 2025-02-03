package controllers

import (
	"bookstore/models"
	"bookstore/repositories"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests made.",
		},
		[]string{"method", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func init() {

	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

type BookController struct {
	Repo *repositories.BookRepository
	Log  *logrus.Logger
}

func NewBookController(repo *repositories.BookRepository) *BookController {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return &BookController{Repo: repo, Log: logger}
}

func (bc *BookController) Metrics(c *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(c.Writer, c.Request)
}

func (bc *BookController) AddBook(c *gin.Context) {

	start := time.Now()

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid JSON when adding book")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		httpRequests.WithLabelValues("POST", "400").Inc()
		httpDuration.WithLabelValues("POST").Observe(time.Since(start).Seconds())
		return
	}

	book, err := bc.Repo.AddBook(book)
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not add book")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add book"})
		httpRequests.WithLabelValues("POST", "500").Inc()
		httpDuration.WithLabelValues("POST").Observe(time.Since(start).Seconds())
		return
	}

	bc.Log.WithFields(logrus.Fields{
		"book_id": book.ID,
		"title":   book.Title,
	}).Info("Book added successfully")

	c.JSON(http.StatusCreated, book)
	httpRequests.WithLabelValues("POST", "201").Inc()
	httpDuration.WithLabelValues("POST").Observe(time.Since(start).Seconds())
}

func (bc *BookController) GetAllBooks(c *gin.Context) {
	// Start timer to track request duration
	start := time.Now()

	books, err := bc.Repo.GetAllBooks()
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Could not retrieve books")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve books"})
		httpRequests.WithLabelValues("GET", "500").Inc()
		httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
		return
	}

	if len(books) == 0 {
		bc.Log.Info("No books available in the bookstore")
		c.JSON(http.StatusOK, gin.H{"message": "No books are available in the bookstore."})
		httpRequests.WithLabelValues("GET", "200").Inc()
		httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
		return
	}

	bc.Log.WithFields(logrus.Fields{
		"book_count": len(books),
	}).Info("Books retrieved successfully")
	c.JSON(http.StatusOK, books)
	httpRequests.WithLabelValues("GET", "200").Inc()
	httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
}

func (bc *BookController) GetBookByID(c *gin.Context) {
	// Start timer to track request duration
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.Param("id"),
		}).Error("Invalid book ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		httpRequests.WithLabelValues("GET", "400").Inc()
		httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
		return
	}

	book, err := bc.Repo.GetBookByID(uint(id))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Book not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		httpRequests.WithLabelValues("GET", "404").Inc()
		httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
		return
	}

	bc.Log.WithFields(logrus.Fields{
		"book_id": book.ID,
		"title":   book.Title,
	}).Info("Book retrieved successfully")
	c.JSON(http.StatusOK, book)
	httpRequests.WithLabelValues("GET", "200").Inc()
	httpDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	// Start timer to track request duration
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.Param("id"),
		}).Error("Invalid book ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		httpRequests.WithLabelValues("DELETE", "400").Inc()
		httpDuration.WithLabelValues("DELETE").Observe(time.Since(start).Seconds())
		return
	}

	_, err = bc.Repo.GetBookByID(uint(id))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Book not found for deletion")
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		httpRequests.WithLabelValues("DELETE", "404").Inc()
		httpDuration.WithLabelValues("DELETE").Observe(time.Since(start).Seconds())
		return
	}

	err = bc.Repo.DeleteBook(uint(id))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Failed to delete book")
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete book"})
		httpRequests.WithLabelValues("DELETE", "500").Inc()
		httpDuration.WithLabelValues("DELETE").Observe(time.Since(start).Seconds())
		return
	}

	bc.Log.WithFields(logrus.Fields{
		"book_id": id,
	}).Info("Book deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
	httpRequests.WithLabelValues("DELETE", "200").Inc()
	httpDuration.WithLabelValues("DELETE").Observe(time.Since(start).Seconds())
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	// Start timer to track request duration
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    c.Param("id"),
		}).Error("Invalid book ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		httpRequests.WithLabelValues("PUT", "400").Inc()
		httpDuration.WithLabelValues("PUT").Observe(time.Since(start).Seconds())
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid JSON when updating book")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		httpRequests.WithLabelValues("PUT", "400").Inc()
		httpDuration.WithLabelValues("PUT").Observe(time.Since(start).Seconds())
		return
	}

	book, err := bc.Repo.UpdateBook(uint(id), updatedBook)
	if err != nil {
		bc.Log.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Book not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		httpRequests.WithLabelValues("PUT", "404").Inc()
		httpDuration.WithLabelValues("PUT").Observe(time.Since(start).Seconds())
		return
	}

	bc.Log.WithFields(logrus.Fields{
		"book_id": book.ID,
		"title":   book.Title,
	}).Info("Book updated successfully")
	c.JSON(http.StatusOK, book)
	httpRequests.WithLabelValues("PUT", "200").Inc()
	httpDuration.WithLabelValues("PUT").Observe(time.Since(start).Seconds())
}
