package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
	"github.com/hungqd/books-service/controller"
)

func GetHandler(bookService book.Service) http.Handler {
	r := gin.Default()

	bookController := controller.NewBookController(bookService)
	r.POST("/books", bookController.CreateBook)

	return r.Handler()
}
