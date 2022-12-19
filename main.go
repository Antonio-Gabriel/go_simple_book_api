package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Clean code", Author: "Robert Cecil Martin", Quantity: 2},
	{ID: "2", Title: "How to build complex system design", Author: "Petter something", Quantity: 3},
	{ID: "3", Title: "Guide to understand machine learning", Author: "Microsoft", Quantity: 2},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookByIdMemory(id string) (*book, error) {
	for i, bk := range books {
		if bk.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookByIdMemory(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookByIdMemory(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var bookModel book

	if err := c.BindJSON(&bookModel); err != nil {
		return
	}

	books = append(books, bookModel)
	c.IndentedJSON(http.StatusCreated, bookModel)
}

func getBookById(c *gin.Context) {
	bookdId := c.Param("id")
	book, err := getBookByIdMemory(bookdId)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("return", returnBook)

	router.Run("localhost:8080")
}
