package book

import "gorm.io/gorm"

type Repository interface {
	SaveBook(book *Book) error
}

type repository struct {
	db *gorm.DB
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
