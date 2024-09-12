package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
	"github.com/jackc/pgx/v5/pgconn"
)

type BookController interface {
	CreateBook(c *gin.Context)
}

type bookController struct {
	service book.Service
}

// CreateBook implements BookController.
func (b *bookController) CreateBook(c *gin.Context) {
	var data book.CreateBook
	c.ShouldBind(&data)
	created, err := b.service.CreateBook(&data)
	if err != nil {
		if err.(*pgconn.PgError) != nil {
			c.JSON(400, gin.H{"error": err})
			return
		} else {
			c.JSON(500, gin.H{"error": err})
			return
		}
	}
	c.JSON(201, created)
}

func NewBookController(service book.Service) BookController {
	return &bookController{service: service}
}
