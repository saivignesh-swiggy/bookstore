package repositories

import (
	"bookstore/metrics"
	"bookstore/models"
	"gorm.io/gorm"
	"time"
)

// SQLiteBookRepository is a concrete implementation of BookRepository using SQLite.
type SQLiteBookRepository struct {
	DB      *gorm.DB
	Metrics *metrics.APIMetrics
}

// NewSQLiteBookRepository creates a new SQLiteBookRepository instance.
func NewSQLiteBookRepository(db *gorm.DB, metrics *metrics.APIMetrics) *SQLiteBookRepository {
	return &SQLiteBookRepository{DB: db, Metrics: metrics}
}

func (r *SQLiteBookRepository) AddBook(book models.Book) (models.Book, error) {
	start := time.Now()
	result := r.DB.Create(&book)
	duration := time.Since(start).Seconds()

	// Track query metrics
	r.Metrics.IncQueryCount("INSERT")
	r.Metrics.MeasureQueryLatency("INSERT", duration)

	return book, result.Error
}

func (r *SQLiteBookRepository) GetAllBooks() ([]models.Book, error) {
	start := time.Now()
	var books []models.Book
	result := r.DB.Find(&books)
	duration := time.Since(start).Seconds()

	// Track query metrics
	r.Metrics.IncQueryCount("SELECT")
	r.Metrics.MeasureQueryLatency("SELECT", duration)

	return books, result.Error
}

func (r *SQLiteBookRepository) GetBookByID(id uint) (models.Book, error) {
	start := time.Now()
	var book models.Book
	result := r.DB.First(&book, id)
	duration := time.Since(start).Seconds()

	// Track query metrics
	r.Metrics.IncQueryCount("SELECT")
	r.Metrics.MeasureQueryLatency("SELECT", duration)

	return book, result.Error
}

func (r *SQLiteBookRepository) DeleteBook(id uint) error {
	start := time.Now()
	err := r.DB.Delete(&models.Book{}, id).Error
	duration := time.Since(start).Seconds()

	// Track query metrics
	r.Metrics.IncQueryCount("DELETE")
	r.Metrics.MeasureQueryLatency("DELETE", duration)

	return err
}

func (r *SQLiteBookRepository) UpdateBook(id uint, updatedBook models.Book) (models.Book, error) {
	start := time.Now()
	var book models.Book
	result := r.DB.First(&book, id)
	if result.Error != nil {
		return book, result.Error
	}

	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Price = updatedBook.Price
	r.DB.Save(&book)

	duration := time.Since(start).Seconds()

	// Track query metrics
	r.Metrics.IncQueryCount("UPDATE")
	r.Metrics.MeasureQueryLatency("UPDATE", duration)

	return book, nil
}
