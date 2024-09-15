package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
	"github.com/hungqd/books-service/controller"
)

func NewHandler(bookService book.Service) http.Handler {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	bookController := controller.NewBookController(bookService)
	bookGroup := r.Group("/books")
	{
		bookGroup.GET("", bookController.GetBooks)
		bookGroup.POST("", bookController.CreateBook)
	}

	return r.Handler()
}
