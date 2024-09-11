package book

type Service interface {
}

type service struct {
	r Repository
}

func NewService(repository Repository) Service {
	return &service{
		r: repository,
	}
}
