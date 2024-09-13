package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
)

type BookController interface {
	GetBooks(c *gin.Context)
	CreateBook(c *gin.Context)
}

type bookController struct {
	service book.Service
}

type bookDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Thumbnail string    `json:"thumbnail"`
	DetailURL string    `json:"detail_url"`
	Title     string    `json:"title"`
	Rating    int       `json:"rating"`
	Price     string    `json:"price"`
	Instock   bool      `json:"instock"`
}

// GetBooks implements BookController.
func (b *bookController) GetBooks(c *gin.Context) {
	books, err := b.service.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bookDtos := make([]bookDto, len(*books))
	for i, book := range *books {
		bookDtos[i] = bookDto{
			ID:        book.ID,
			CreatedAt: book.CreatedAt,
			Thumbnail: book.Thumbnail,
			DetailURL: book.DetailURL,
			Title:     book.Title,
			Rating:    book.Rating,
			Price:     book.Price,
			Instock:   book.Instock,
		}
	}
	c.JSON(http.StatusOK, bookDtos)
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
		switch err {
		case book.ErrBookAlreadyExists:
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
