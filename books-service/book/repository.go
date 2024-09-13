package book

import "gorm.io/gorm"

type Repository interface {
	GetBooks() (*[]Book, error)
	SaveBook(book *Book) error
}

type repository struct {
	db *gorm.DB
}

// GetBooks implements Repository.
func (r *repository) GetBooks() (*[]Book, error) {
	var books []Book
	result := r.db.Find(&books)
	return &books, result.Error
}

// SaveBook implements Repository.
func (r *repository) SaveBook(book *Book) error {
	result := r.db.Create(book)
	return result.Error
}

func NewRepository(db *gorm.DB) Repository {
	// db.AutoMigrate(&Book{})
	return &repository{
		db,
	}
}
