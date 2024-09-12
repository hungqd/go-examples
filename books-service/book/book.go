package book

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Thumbnail string
	DetailURL string `gorm:"unique"`
	Title     string
	Rating    int
	Price     string
	Instock   bool
}
