package routes

import (
	"bookstore/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(bookController *controllers.BookController) *gin.Engine {
	r := gin.Default()
	
	r.POST("/books", bookController.AddBook)
	r.GET("/books", bookController.GetAllBooks)
	r.GET("/books/:id", bookController.GetBookByID)
	r.PUT("/books/:id", bookController.UpdateBook)
	r.DELETE("/books/:id", bookController.DeleteBook)
	return r
}
