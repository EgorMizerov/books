package models

import (
	"errors"
	"github.com/google/uuid"
)

type Book struct {
	ID     uuid.UUID
	Title  string
	Author Author
}

func NewBook(title string, bookId uuid.UUID, author Author) (Book, error) {
	if title == "" {
		return Book{}, errors.New("book title must not be empty")
	}
	return Book{
		ID:     bookId,
		Title:  title,
		Author: author,
	}, nil
}
