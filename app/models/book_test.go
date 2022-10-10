package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBookErrorIfInvalidTitle(t *testing.T) {
	result, err := NewBook("", uuid.New(), Author{})

	assert.EqualError(t, err, "book title must not be empty")
	assert.Equal(t, Book{}, result)
}

func TestNewBook(t *testing.T) {
	bookId := uuid.New()
	bookTitle := "test_title"
	author := Author{
		ID:   uuid.New(),
		Name: "test_name",
	}

	result, err := NewBook(bookTitle, bookId, author)

	assert.NoError(t, err)
	assert.Equal(t, Book{
		ID:     bookId,
		Title:  bookTitle,
		Author: author,
	}, result)
}
