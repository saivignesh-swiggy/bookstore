package repositories

import "bookstore/models"

type BookRepository interface {
	AddBook(book models.Book) (models.Book, error)
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id uint) (models.Book, error)
	UpdateBook(id uint, updatedBook models.Book) (models.Book, error)
	DeleteBook(id uint) error
}
