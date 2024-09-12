package book

type CreateBook struct {
	Thumbnail string
	DetailURL string
	Title     string
	Rating    int
	Price     string
	Instock   bool
}

type Service interface {
	CreateBook(book *CreateBook) (*Book, error)
}

type service struct {
	r Repository
}

// CreateBook implements Service.
func (s *service) CreateBook(book *CreateBook) (*Book, error) {
	entity := Book{
		Thumbnail: book.Thumbnail,
		DetailURL: book.DetailURL,
		Title:     book.Title,
		Rating:    book.Rating,
		Price:     book.Price,
		Instock:   book.Instock,
	}
	err := s.r.SaveBook(&entity)
	return &entity, err
}

func NewService(repository Repository) Service {
	return &service{
		r: repository,
	}
}
