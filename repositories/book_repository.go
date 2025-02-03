package repositories

import (
	"bookstore/db"
	"bookstore/models"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

func (r *BookRepository) AddBook(book models.Book) (models.Book, error) {
	result := db.DB.Create(&book)
	return book, result.Error
}

func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	result := db.DB.Find(&books)
	return books, result.Error
}

func (r *BookRepository) GetBookByID(id uint) (models.Book, error) {
	var book models.Book
	result := db.DB.First(&book, id)
	return book, result.Error
}

func (r *BookRepository) DeleteBook(id uint) error {
	return db.DB.Delete(&models.Book{}, id).Error

}

func (r *BookRepository) UpdateBook(id uint, updatedBook models.Book) (models.Book, error) {
	var book models.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		return book, result.Error
	}
	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Price = updatedBook.Price
	db.DB.Save(&book)
	return book, nil
}
