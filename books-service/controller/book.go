package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
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
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := b.service.CreateBook(&data)
	if err != nil {
		switch err.(type) {
		case *pgconn.PgError, sqlite3.Error:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(201, gin.H{
		"id":         created.ID,
		"created_at": created.CreatedAt,
		"thumbnail":  created.Thumbnail,
		"detail_url": created.DetailURL,
		"title":      created.Title,
		"rating":     created.Rating,
		"price":      created.Price,
		"instock":    created.Instock,
	})
}

func NewBookController(service book.Service) BookController {
	return &bookController{service: service}
}
