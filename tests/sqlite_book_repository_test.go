package repositories_test

import (
	"bookstore/metrics"
	"bookstore/models"
	"bookstore/repositories" // Import the repositories package
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// Setup an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to connect to in-memory database")
	assert.NoError(t, db.AutoMigrate(&models.Book{}), "Failed to migrate database")
	return db
}

func TestAddBook(t *testing.T) {
	db := setupTestDB(t)
	mockMetrics := metrics.NewAPIMetrics()
	repo := repositories.NewSQLiteBookRepository(db, mockMetrics) // FIXED

	book := models.Book{Title: "Test Book", Author: "Author X", Price: 10.99}
	addedBook, err := repo.AddBook(book)

	assert.NoError(t, err, "Failed to add book")
	assert.NotZero(t, addedBook.ID, "Book ID should be set after insertion")
	assert.Equal(t, book.Title, addedBook.Title)
}

func TestGetAllBooks(t *testing.T) {
	db := setupTestDB(t)
	mockMetrics := metrics.NewAPIMetrics()
	repo := repositories.NewSQLiteBookRepository(db, mockMetrics) // FIXED

	// Insert test data
	repo.AddBook(models.Book{Title: "Book 1", Author: "Author A", Price: 15.99})
	repo.AddBook(models.Book{Title: "Book 2", Author: "Author B", Price: 20.99})

	books, err := repo.GetAllBooks()

	assert.NoError(t, err)
	assert.Len(t, books, 2, "Expected 2 books in database")
}

func TestGetBookByID(t *testing.T) {
	db := setupTestDB(t)
	mockMetrics := metrics.NewAPIMetrics()
	repo := repositories.NewSQLiteBookRepository(db, mockMetrics) // FIXED

	book := models.Book{Title: "Unique Book", Author: "Unique Author", Price: 9.99}
	addedBook, _ := repo.AddBook(book)

	foundBook, err := repo.GetBookByID(addedBook.ID)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, foundBook.Title, "Fetched book title should match")
}

func TestUpdateBook(t *testing.T) {
	db := setupTestDB(t)
	mockMetrics := metrics.NewAPIMetrics()
	repo := repositories.NewSQLiteBookRepository(db, mockMetrics) // FIXED

	// Add a book
	book := models.Book{Title: "Old Title", Author: "Author Y", Price: 12.99}
	addedBook, _ := repo.AddBook(book)

	// Update it
	updatedBook := models.Book{Title: "New Title", Author: "Author Y", Price: 14.99}
	updated, err := repo.UpdateBook(addedBook.ID, updatedBook)

	assert.NoError(t, err)
	assert.Equal(t, "New Title", updated.Title, "Title should be updated")
	assert.Equal(t, 14.99, updated.Price, "Price should be updated")
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	mockMetrics := metrics.NewAPIMetrics()
	repo := repositories.NewSQLiteBookRepository(db, mockMetrics) // FIXED

	book := models.Book{Title: "Delete Me", Author: "Author Z", Price: 18.99}
	addedBook, _ := repo.AddBook(book)

	err := repo.DeleteBook(addedBook.ID)
	assert.NoError(t, err, "Book should be deleted successfully")

	// Try to get deleted book
	_, err = repo.GetBookByID(addedBook.ID)
	assert.Error(t, err, "Deleted book should not be found")
}
