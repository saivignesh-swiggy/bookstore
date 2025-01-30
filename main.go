package main

import (
	"fmt"
	"strconv"

	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
/*type Book struct {
	id     int
	title  string
	author string
	price  float64
}*/
type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	//fmt.Println("hello world")
	r := gin.Default()
	var books = make(map[int]Book)
	var titleAuthorIndex = make(map[string]bool)
	var nextid = 1
	// adds a book
	r.POST("/books", func(c *gin.Context) {
		var book Book
		err := c.BindJSON(&book)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		/*for _, existingBook := range books {
			if existingBook.Title == book.Title && existingBook.Author == book.Author {
				c.JSON(http.StatusConflict, gin.H{"error": "Book already present"})
				return
			}
		}*/
		key := fmt.Sprintf("%s:%s", book.Title, book.Author)
		if _, exists := titleAuthorIndex[key]; exists {
			c.JSON(http.StatusConflict, gin.H{"message": "Book already present"})
			return
		}
		book.ID = nextid
		books[nextid] = book
		titleAuthorIndex[key] = true
		nextid++
		c.JSON(http.StatusCreated, book)
	})
	// gets all books
	r.GET("/books", func(c *gin.Context) {
		if len(books) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No books are available in the bookstore."})
			return
		}
		var allbooks []Book
		for _, book := range books {
			allbooks = append(allbooks, book)
		}

		c.JSON(http.StatusOK, allbooks)
	})
	// get a book by id
	r.GET("/books/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		book, exists := books[id]
		if exists {
			c.JSON(http.StatusOK, book)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		}
	})
	r.DELETE("/books/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		if book, exists := books[id]; exists {
			key := fmt.Sprintf("%s:%s", book.Title, book.Author)
			delete(titleAuthorIndex, key)
			delete(books, id)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		}
	})
	// updates book
	r.PUT("/books/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		var updatedBook Book
		if err := c.BindJSON(&updatedBook); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		book, exists := books[id]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		oldkey := fmt.Sprintf("%s:%s", book.Title, book.Author)
		delete(titleAuthorIndex, oldkey)

		newkey := fmt.Sprintf("%s:%s", updatedBook.Title, updatedBook.Author)
		titleAuthorIndex[newkey] = true

		book.Title = updatedBook.Title
		book.Author = updatedBook.Author
		book.Price = updatedBook.Price

		books[id] = book

		c.JSON(http.StatusOK, book)
	})
	r.Run(":8080")

}

/*
curl -X GET "http://localhost:8080/books/1"

curl -X GET "http://localhost:8080/books"

curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Go Programming",
           "author": "John Doe",
           "price": 29.99
         }'
curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{
           "title": "The White Tiger",
           "author": "Aravind Adiga",
           "price": 14.99
         }'
curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Whispers of the Himalayas",
           "author": "Aditi Sharma",
           "price": 19.99
         }'
curl -X PUT http://localhost:8080/books/3 \
     -H "Content-Type: application/json" \
     -d '{
           "title": "The White Tiger",
           "author": "Aravind Adiga",
           "price": 17.99
         }'

curl -X DELETE "http://localhost:8080/books/1"

*/
