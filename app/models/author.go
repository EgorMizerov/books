package models

import (
	"errors"
	"github.com/google/uuid"
)

type Author struct {
	ID   uuid.UUID
	Name string
}

func NewAuthor(name string, authorId uuid.UUID) (Author, error) {
	if name == "" {
		return Author{}, errors.New("author name must not be empty")
	}
	return Author{
		ID:   authorId,
		Name: name,
	}, nil
}
