package book

import "gorm.io/gorm"

type Repository interface {
	ExistByDetailURL(detailURL string) (bool, error)
	GetBooks() (*[]Book, error)
	SaveBook(book *Book) error
}

type repository struct {
	db *gorm.DB
}

// ExistByDetailURL implements Repository.
func (r *repository) ExistByDetailURL(detailURL string) (bool, error) {
	var count int64
	result := r.db.Model(&Book{}).Where("detail_url = ?", detailURL).Count(&count)
	return count > 0, result.Error
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
