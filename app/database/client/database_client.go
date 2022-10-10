package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/egormizerov/books/app/models"
)

// Query template to create author.
var createAuthorQuery = `INSERT INTO authors (id, name) VALUES (:id, :name)`

// Query template to create book.
var createBookQuery = `INSERT INTO books (id, title, author_id) VALUES (:id, :title, :author_id)`

// Query template to get book by id.
var getBookByIdQuery = `SELECT id, title, author_id FROM books WHERE id=:book_id`

// Query template to get author by id.
var getAuthorByIdQuery = `SELECT id, name FROM authors WHERE id=:author_id`

// Query template to get book's by author id.
var getBooksByAuthorIdQuery = `SELECT id, title, author_id FROM books WHERE author_id=:author_id`

type DatabaseClient struct {
	db *sqlx.DB
}

func NewDatabaseClient(db *sqlx.DB) *DatabaseClient {
	return &DatabaseClient{db: db}
}

type createAuthorArguments struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (self *DatabaseClient) CreateAuthor(ctx context.Context, author models.Author) error {
	_, err := self.db.NamedExecContext(ctx, createAuthorQuery, createAuthorArguments{
		ID:   author.ID,
		Name: author.Name,
	})
	return err
}

type createBookArguments struct {
	ID       uuid.UUID `db:"id"`
	Title    string    `db:"title"`
	AuthorId uuid.UUID `db:"author_id"`
}

func (self *DatabaseClient) CreateBook(ctx context.Context, book models.Book) error {
	_, err := self.db.NamedExecContext(ctx, createBookQuery, createBookArguments{
		ID:       book.ID,
		Title:    book.Title,
		AuthorId: book.Author.ID,
	})
	return err
}

type getBookArguments struct {
	BookId uuid.UUID `db:"book_id"`
}

func (self *DatabaseClient) GetBookById(ctx context.Context, bookId uuid.UUID) (models.Book, error) {
	rows, err := self.db.NamedQueryContext(ctx, getBookByIdQuery, getBookArguments{
		BookId: bookId,
	})
	if err != nil {
		return models.Book{}, err
	}
	if !rows.Next() {
		return models.Book{}, errors.New("no rows in result set")
	}
	book := models.Book{Author: models.Author{}}
	err = rows.Scan(&book.ID, &book.Title, &book.Author.ID)
	if err != nil {
		return models.Book{}, fmt.Errorf("failed to scan row: %s", err)
	}

	return book, nil
}

type getAuthorArguments struct {
	AuthorId uuid.UUID `db:"author_id"`
}

func (self *DatabaseClient) GetAuthorById(ctx context.Context, authorId uuid.UUID) (models.Author, error) {
	rows, err := self.db.NamedQueryContext(ctx, getAuthorByIdQuery, getAuthorArguments{
		AuthorId: authorId,
	})
	if err != nil {
		return models.Author{}, err
	}
	if !rows.Next() {
		return models.Author{}, errors.New("no rows in result set")
	}
	var author models.Author
	err = rows.Scan(&author.ID, &author.Name)
	if err != nil {
		return models.Author{}, fmt.Errorf("failed to scan row: %s", err)
	}

	return author, nil
}

type getBooksByAuthorIdArguments struct {
	AuthorId uuid.UUID `db:"author_id"`
}

func (self *DatabaseClient) GetBooksByAuthorId(ctx context.Context, authorId uuid.UUID) ([]models.Book, error) {
	rows, err := self.db.NamedQueryContext(ctx, getBooksByAuthorIdQuery, getBooksByAuthorIdArguments{
		AuthorId: authorId,
	})
	if err != nil {
		return nil, err
	}

	var books []models.Book
	for rows.Next() {
		book := models.Book{Author: models.Author{}}
		err = rows.Scan(&book.ID, &book.Title, &book.Author.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %s", err)
		}
		books = append(books, book)
	}

	return books, nil
}
