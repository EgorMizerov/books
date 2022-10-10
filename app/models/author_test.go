package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorErrorIfInvalidName(t *testing.T) {
	result, err := NewAuthor("", uuid.New())

	assert.EqualError(t, err, "author name must not be empty")
	assert.Equal(t, Author{}, result)
}

func TestNewAuthor(t *testing.T) {
	authorId := uuid.New()
	authorName := "test_name"

	result, err := NewAuthor(authorName, authorId)

	assert.NoError(t, err)
	assert.Equal(t, Author{
		ID:   authorId,
		Name: authorName,
	}, result)
}
