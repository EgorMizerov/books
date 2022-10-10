package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/egormizerov/books/app/models"
	logcontext "github.com/egormizerov/books/pkg/log/context"
	"github.com/egormizerov/books/pkg/wrappers"
)

type DatabaseClient interface {
	CreateAuthor(ctx context.Context, author models.Author) error
	CreateBook(ctx context.Context, book models.Book) error
	GetBookById(ctx context.Context, bookId uuid.UUID) (models.Book, error)
	GetAuthorById(ctx context.Context, authorId uuid.UUID) (models.Author, error)
	GetBooksByAuthorId(ctx context.Context, authorId uuid.UUID) ([]models.Book, error)
}

type Service struct {
	DatabaseClient DatabaseClient
	uuid           wrappers.UUIDWrapper
}

func NewService(databaseClient DatabaseClient, uuid wrappers.UUIDWrapper) *Service {
	return &Service{
		DatabaseClient: databaseClient,
		uuid:           uuid,
	}
}

func (self *Service) CreateBook(ctx context.Context, bookTitle string, authorId uuid.UUID) error {
	book, err := models.NewBook(bookTitle, self.uuid.New(), models.Author{ID: authorId})
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("author_id", authorId.String()).
			WithError(err).
			Error("failed to init book")
		return fmt.Errorf("failed to init book: %s", err)
	}

	err = self.DatabaseClient.CreateBook(ctx, book)
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("author_id", authorId.String()).
			WithField("book_title", bookTitle).
			WithError(err).
			Error("failed to create book")
		return fmt.Errorf("failed to create book: %s", err)
	}

	return nil
}

func (self *Service) CreateAuthor(ctx context.Context, authorName string) error {
	author, err := models.NewAuthor(authorName, self.uuid.New())
	if err != nil {
		logcontext.FromContext(ctx).
			WithError(err).
			Error("failed to init author")
		return fmt.Errorf("failed to init author: %s", err)
	}

	err = self.DatabaseClient.CreateAuthor(ctx, author)
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("author_name", authorName).
			WithError(err).
			Error("failed to create author")
		return fmt.Errorf("failed to create author: %s", err)
	}

	return nil
}

func (self *Service) GetBook(ctx context.Context, bookId uuid.UUID) (models.Book, error) {
	book, err := self.DatabaseClient.GetBookById(ctx, bookId)
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("book_id", bookId.String()).
			WithError(err).
			Error("failed to get book")
		return models.Book{}, fmt.Errorf("failed to get book by id: %s", err)
	}
	bookAuthor, err := self.DatabaseClient.GetAuthorById(ctx, book.Author.ID)
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("book_id", bookId.String()).
			WithError(err).
			Error("failed to get author")
		return models.Book{}, fmt.Errorf("failed to get author by id: %s", err)
	}
	book.Author = bookAuthor

	return book, nil
}

func (self *Service) GetAuthorsBooks(ctx context.Context, authorId uuid.UUID) ([]models.Book, error) {
	books, err := self.DatabaseClient.GetBooksByAuthorId(ctx, authorId)
	if err != nil {
		logcontext.FromContext(ctx).
			WithField("author_id", authorId.String()).
			WithError(err).
			Error("failed to get books by author id")
		return nil, fmt.Errorf("failed to get books by author id: %s", err)
	}

	return books, nil
}
