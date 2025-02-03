package tests

import (
	"bookstore/controllers"
	"bookstore/db"
	"bookstore/models"
	"bookstore/repositories"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupTestDB() {
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to initialize test database")
	}
	db.DB.AutoMigrate(&models.Book{})
}

func setupTestController() (*controllers.BookController, *gin.Engine) {
	repo := repositories.NewBookRepository()
	controller := controllers.NewBookController(repo)

	// Create a test router
	router := gin.Default()
	router.POST("/books", controller.AddBook)
	router.GET("/books", controller.GetAllBooks)
	router.GET("/books/:id", controller.GetBookByID)
	router.PUT("/books/:id", controller.UpdateBook)
	router.DELETE("/books/:id", controller.DeleteBook)

	return controller, router
}

func TestBookController(t *testing.T) {
	setupTestDB()
	_, router := setupTestController()

	t.Run("AddBook", func(t *testing.T) {
		book := models.Book{Title: "Test Book", Author: "Test Author", Price: 9.99}
		body, _ := json.Marshal(book)

		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetAllBooks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetBookByID - Existing", func(t *testing.T) {
		book := models.Book{Title: "Existing Book", Author: "Author", Price: 15.99}
		db.DB.Create(&book) // Manually add a book

		req, _ := http.NewRequest("GET", "/books/"+strconv.Itoa(int(book.ID)), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetBookByID - Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("UpdateBook - Existing", func(t *testing.T) {
		book := models.Book{Title: "Old Book", Author: "Old Author", Price: 20.00}
		db.DB.Create(&book)

		updatedBook := models.Book{Title: "Updated Book", Author: "Updated Author", Price: 25.00}
		body, _ := json.Marshal(updatedBook)

		req, _ := http.NewRequest("PUT", "/books/"+strconv.Itoa(int(book.ID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateBook - Not Found", func(t *testing.T) {
		updatedBook := models.Book{Title: "Non-existent", Author: "Author", Price: 10.00}
		body, _ := json.Marshal(updatedBook)

		req, _ := http.NewRequest("PUT", "/books/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DeleteBook - Existing", func(t *testing.T) {
		book := models.Book{Title: "To Be Deleted", Author: "Author", Price: 5.99}
		db.DB.Create(&book)

		req, _ := http.NewRequest("DELETE", "/books/"+strconv.Itoa(int(book.ID)), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteBook - Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/books/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
