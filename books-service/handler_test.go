package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hungqd/books-service/book"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MustMarshal(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.DebugMode)

	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&book.Book{})
	defer os.Remove("test.db")

	repo := book.NewRepository(db)
	service := book.NewService(repo)

	handler := GetHandler(service)

	tests := []struct {
		name             string
		body             string
		expectStatusCode int
		assertFunctions  [](func(*http.Response, *http.Request) error)
	}{
		{
			name: "Create book successfully",
			body: `{
				"thumbnail": "http://example.com/thumbnail.jpg",
				"detailUrl": "http://example.com/detail",
				"title": "Example Book Title",
				"rating": 5,
				"price": "19.99",
				"instock": true
			}`,
			expectStatusCode: http.StatusCreated,
			assertFunctions: [](func(*http.Response, *http.Request) error){
				jsonpath.Equal("$.title", "Example Book Title"),
				jsonpath.Present("$.rating"),
				jsonpath.Equal("$.price", "19.99"),
				jsonpath.Equal("$.instock", true),
				jsonpath.Present("$.created_at"),
				jsonpath.Present("$.id"),
			},
		},
		{
			name: "Create book failed",
			body: `{
				"thumbnail": "http://example.com/thumbnail.jpg",
				"detailUrl": "http://example.com/detail", // Duplicate detailUrl
				"title": "Example Book Title",
				"rating": 5,
				"price": "19.99",
				"instock": true
			}`,
			expectStatusCode: http.StatusBadRequest,
			assertFunctions: [](func(*http.Response, *http.Request) error){
				jsonpath.Present("$.error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(T *testing.T) {
			at := apitest.New().
				Handler(handler).
				Post("/books").
				JSON(tt.body).
				Expect(t)
			for _, fn := range tt.assertFunctions {
				at.Assert(fn)
			}
			at.Status(tt.expectStatusCode).
				End()
		})
	}
}
