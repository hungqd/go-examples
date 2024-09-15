package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service interface {
	SaveBook(book *Book) error
}

type service struct{}

// SaveBook implements Service.
func (s *service) SaveBook(book *Book) error {

	json, err := json.Marshal(book)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"http://localhost:8080/books",
		"application/json",
		bytes.NewBuffer(json),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("save book: %v", resp.StatusCode)
	}
	return nil
}

func NewService() Service {
	return &service{}
}
