package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hungqd/books-service/book"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateBook(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&book.Book{})
	repo := book.NewRepository(db)
	service := book.NewService(repo)

	handler := GetHandler(service)

	tests := []struct {
		name             string
		body             string
		expectStatusCode int
		expectRespBody   string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(T *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/books", strings.NewReader(tt.body))
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatusCode, w.Code)
			if tt.expectRespBody != "" {
				assert.Equal(t, tt.expectRespBody, w.Body.Bytes())
			}
		})
	}
}
