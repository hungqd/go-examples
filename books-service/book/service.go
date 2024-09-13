package book

type CreateBook struct {
	Thumbnail string `binding:"required"`
	DetailURL string `binding:"required"`
	Title     string `binding:"required"`
	Rating    int    `binding:"required"`
	Price     string `binding:"required"`
	Instock   bool   `binding:"required"`
}

type Service interface {
	GetBooks() (*[]Book, error)
	CreateBook(book *CreateBook) (*Book, error)
}

type service struct {
	r Repository
}

// GetBooks implements Service.
func (s *service) GetBooks() (*[]Book, error) {
	books, err := s.r.GetBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
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
