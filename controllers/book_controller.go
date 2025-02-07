package controllers

import (
	"bookstore/metrics"
	"bookstore/models"
	"bookstore/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type BookController struct {
	Repo    repositories.BookRepository
	Log     *logrus.Logger
	Metrics *metrics.APIMetrics
}

func (bc *BookController) handleSuccess(c *gin.Context, action string, statusCode int, data interface{}, fields logrus.Fields, start time.Time) {
	bc.Log.WithFields(fields).Info("Request successful")
	c.JSON(statusCode, data)
	bc.Metrics.IncRequestCount(action, fmt.Sprintf("%d", statusCode))
	bc.Metrics.MeasureRequestLatency(action, time.Since(start).Seconds())
}
func (bc *BookController) handleInfo(c *gin.Context, action string, statusCode int, msg string, start time.Time) {
	bc.Log.Info(msg)
	c.JSON(statusCode, gin.H{"message": msg})
	bc.Metrics.IncRequestCount(action, fmt.Sprintf("%d", statusCode))
	bc.Metrics.MeasureRequestLatency(action, time.Since(start).Seconds())
}

func (bc *BookController) handleErrorResponse(c *gin.Context, statusCode int, errMsg, logMsg, action string, err error, start time.Time) {
	bc.Log.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Error(logMsg)

	c.JSON(statusCode, gin.H{"error": errMsg})
	bc.Metrics.IncRequestCount(action, fmt.Sprintf("%d", statusCode))
	bc.Metrics.MeasureRequestLatency(action, time.Since(start).Seconds())
}

func NewBookController(repo repositories.BookRepository, metrics *metrics.APIMetrics) *BookController {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return &BookController{Repo: repo, Log: logger, Metrics: metrics}
}

func (bc *BookController) AddBook(c *gin.Context) {
	start := time.Now()

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		/*bc.Log.WithFields(logrus.Fields{"error": err.Error()}).Error("Invalid JSON when adding book")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		bc.Metrics.IncRequestCount("POST", "400")
		bc.Metrics.MeasureRequestLatency("POST", time.Since(start).Seconds())*/
		bc.handleErrorResponse(c, http.StatusBadRequest, "Invalid JSON", "Invalid JSON when adding book", "addBook", err, start)

		return
	}

	book, err := bc.Repo.AddBook(book)
	if err != nil {
		/*bc.Log.WithFields(logrus.Fields{"error": err.Error()}).Error("Could not add book")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add book"})
		bc.Metrics.IncRequestCount("POST", "500")
		bc.Metrics.MeasureRequestLatency("POST", time.Since(start).Seconds())*/
		bc.handleErrorResponse(c, http.StatusInternalServerError, "Could not add book", "Could not add book", "addBook", err, start)

		return
	}
	bc.handleSuccess(c, "addBook", http.StatusCreated, book, logrus.Fields{"book_id": book.ID, "title": book.Title}, start)
	bc.Metrics.IncBooksAdded("success")
}
func (bc *BookController) GetAllBooks(c *gin.Context) {
	start := time.Now()

	books, err := bc.Repo.GetAllBooks()
	if err != nil {
		bc.handleErrorResponse(c, http.StatusInternalServerError, "Could not retrieve books", "Could not retrieve books", "getAllBooks", err, start)
		return
	}

	if len(books) == 0 {
		bc.handleInfo(c, "getAllBooks", http.StatusOK, "No books are available in the bookstore.", start)
		return
	}

	bookCount := len(books)
	bc.Metrics.UpdateBooksAvailable(bookCount)
	bc.handleSuccess(c, "getAllBooks", http.StatusOK, books, logrus.Fields{"book_count": bookCount}, start)
}

func (bc *BookController) GetBookByID(c *gin.Context) {
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.handleErrorResponse(c, http.StatusBadRequest, "Invalid book ID", "Invalid book ID", "getBook", err, start)
		return
	}

	book, err := bc.Repo.GetBookByID(uint(id))
	if err != nil {
		bc.handleErrorResponse(c, http.StatusNotFound, "Book not found", "Book not found", "getBook", err, start)
		return
	}

	bc.handleSuccess(c, "getBook", http.StatusOK, book, logrus.Fields{"book_id": book.ID, "title": book.Title}, start)

}
func (bc *BookController) DeleteBook(c *gin.Context) {
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.handleErrorResponse(c, http.StatusBadRequest, "Invalid book ID", "Invalid book ID", "deleteBook", err, start)
		return
	}

	_, err = bc.Repo.GetBookByID(uint(id))
	if err != nil {
		bc.handleErrorResponse(c, http.StatusNotFound, "Book not found", "Book not found for deletion", "deleteBook", err, start)
		return
	}

	if err = bc.Repo.DeleteBook(uint(id)); err != nil {
		bc.handleErrorResponse(c, http.StatusInternalServerError, "Failed to delete book", "Failed to delete book", "deleteBook", err, start)
		return
	}

	bc.handleSuccess(c, "deleteBook", http.StatusOK, gin.H{"message": "Book deleted"}, logrus.Fields{"book_id": id}, start)
}
func (bc *BookController) UpdateBook(c *gin.Context) {
	start := time.Now()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		bc.handleErrorResponse(c, http.StatusBadRequest, "Invalid book ID", "Invalid book ID", "updateBook", err, start)
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		bc.handleErrorResponse(c, http.StatusBadRequest, "Invalid JSON", "Invalid JSON when updating book", "updateBook", err, start)
		return
	}

	book, err := bc.Repo.UpdateBook(uint(id), updatedBook)
	if err != nil {
		bc.handleErrorResponse(c, http.StatusNotFound, "Book not found", "Book not found for update", "updateBook", err, start)
		return
	}

	bc.handleSuccess(c, "updateBook", http.StatusOK, book, logrus.Fields{"book_id": book.ID, "title": book.Title}, start)
}
