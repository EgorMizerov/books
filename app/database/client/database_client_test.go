package client

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"database/sql/driver"
	"github.com/egormizerov/books/app/models"
)

var (
	createAuthorQueryMatcher       = regexp.QuoteMeta(`INSERT INTO authors (id, name) VALUES (?, ?)`)
	createBookQueryMatcher         = regexp.QuoteMeta(`INSERT INTO books (id, title, author_id) VALUES (?, ?, ?)`)
	getBookByIdQueryMatcher        = `SELECT id, title, author_id FROM books WHERE id=?`
	getAuthorByIdQueryMatcher      = `SELECT id, name FROM authors WHERE id=?`
	getBooksByAuthorIdQueryMatcher = `SELECT id, title, author_id FROM books WHERE author_id=?`
)

type DatabaseClientTests struct {
	suite.Suite
	client  DatabaseClient
	context context.Context
	sqlxDB  *sqlx.DB
	sqlMock sqlmock.Sqlmock

	testError error
	book      models.Book
	author    models.Author
}

func TestDatabaseClient(t *testing.T) {
	suite.Run(t, new(DatabaseClientTests))
}

func (self *DatabaseClientTests) SetupTest() {
	mockDatabaseConnection, sqlMock, err := sqlmock.New()
	self.Require().Nil(err)

	self.sqlxDB = sqlx.NewDb(mockDatabaseConnection, "sqlmock")
	self.sqlMock = sqlMock
	self.client = DatabaseClient{db: self.sqlxDB}
	self.context = context.Background()

	self.author = models.Author{
		ID:   uuid.New(),
		Name: "test_author",
	}
	self.book = models.Book{
		ID:     uuid.New(),
		Title:  "test_title",
		Author: models.Author{ID: self.author.ID},
	}
	self.testError = errors.New("test_error")
}

func (self *DatabaseClientTests) TestNewDatabaseClient() {
	result := NewDatabaseClient(self.sqlxDB)

	self.Equal(&DatabaseClient{
		db: self.sqlxDB,
	}, result)
}

func (self *DatabaseClientTests) TestCreateAuthorErrorIfSqlExecFailed() {
	self.sqlMock.
		ExpectExec(createAuthorQueryMatcher).
		WithArgs(self.author.ID, self.author.Name).
		WillReturnError(self.testError)

	err := self.client.CreateAuthor(self.context, self.author)

	self.EqualError(err, self.testError.Error())
}

func (self *DatabaseClientTests) TestCreateAuthor() {
	self.sqlMock.
		ExpectExec(createAuthorQueryMatcher).
		WithArgs(self.author.ID, self.author.Name).
		WillReturnResult(driver.ResultNoRows)

	err := self.client.CreateAuthor(self.context, self.author)

	self.NoError(err)
}

func (self *DatabaseClientTests) TestCreateBookErrorIfSqlExecFailed() {
	self.sqlMock.
		ExpectExec(createBookQueryMatcher).
		WithArgs(self.book.ID, self.book.Title, self.book.Author.ID).
		WillReturnError(self.testError)

	err := self.client.CreateBook(self.context, self.book)

	self.EqualError(err, self.testError.Error())
}

func (self *DatabaseClientTests) TestCreateBook() {
	self.sqlMock.
		ExpectExec(createBookQueryMatcher).
		WithArgs(self.book.ID, self.book.Title, self.book.Author.ID).
		WillReturnResult(driver.ResultNoRows)

	err := self.client.CreateBook(self.context, self.book)

	self.NoError(err)
}

func (self *DatabaseClientTests) TestGetBookByIdErrorIfSqlQueryFailed() {
	self.sqlMock.
		ExpectQuery(getBookByIdQueryMatcher).
		WithArgs(self.book.ID).
		WillReturnError(self.testError)

	result, err := self.client.GetBookById(self.context, self.book.ID)

	self.EqualError(err, self.testError.Error())
	self.Equal(models.Book{}, result)
}

func (self *DatabaseClientTests) TestGetBookByIdErrorIfNoRows() {
	self.sqlMock.
		ExpectQuery(getBookByIdQueryMatcher).
		WithArgs(self.book.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}))

	result, err := self.client.GetBookById(self.context, self.book.ID)

	self.EqualError(err, "no rows in result set")
	self.Equal(models.Book{}, result)
}

func (self *DatabaseClientTests) TestGetBookByIdErrorIfScanRowFailed() {
	rows := sqlmock.NewRows([]string{"not_book_field"}).AddRow(true)
	self.sqlMock.
		ExpectQuery(getBookByIdQueryMatcher).
		WithArgs(self.book.ID).
		WillReturnRows(rows)

	result, err := self.client.GetBookById(self.context, self.book.ID)

	self.ErrorContains(err, "failed to scan row")
	self.Equal(models.Book{}, result)
}

func (self *DatabaseClientTests) TestGetBookById() {
	rows := sqlmock.NewRows([]string{"id", "title", "author_id"}).
		AddRow(self.book.ID, self.book.Title, self.book.Author.ID)
	self.sqlMock.
		ExpectQuery(getBookByIdQueryMatcher).
		WithArgs(self.book.ID).
		WillReturnRows(rows)

	result, err := self.client.GetBookById(self.context, self.book.ID)

	self.NoError(err)
	self.Equal(self.book, result)
}

func (self *DatabaseClientTests) TestGetAuthorByIdErrorIfSqlQueryFailed() {
	self.sqlMock.
		ExpectQuery(getAuthorByIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnError(self.testError)

	result, err := self.client.GetAuthorById(self.context, self.author.ID)

	self.EqualError(err, self.testError.Error())
	self.Equal(models.Author{}, result)
}

func (self *DatabaseClientTests) TestGetAuthorByIdErrorIfNoRows() {
	self.sqlMock.
		ExpectQuery(getAuthorByIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	result, err := self.client.GetAuthorById(self.context, self.author.ID)

	self.EqualError(err, "no rows in result set")
	self.Equal(models.Author{}, result)
}

func (self *DatabaseClientTests) TestGetAuthorByIdErrorIfScanRowFailed() {
	rows := sqlmock.NewRows([]string{"not_author_field"}).AddRow(true)
	self.sqlMock.
		ExpectQuery(getAuthorByIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(rows)

	result, err := self.client.GetAuthorById(self.context, self.author.ID)

	self.ErrorContains(err, "failed to scan row")
	self.Equal(models.Author{}, result)
}

func (self *DatabaseClientTests) TestGetAuthorById() {
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(self.author.ID, self.author.Name)
	self.sqlMock.
		ExpectQuery(getAuthorByIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(rows)

	result, err := self.client.GetAuthorById(self.context, self.author.ID)

	self.NoError(err)
	self.Equal(self.author, result)
}

func (self *DatabaseClientTests) TestGetBooksByAuthorIdErrorIfSqlQueryFailed() {
	self.sqlMock.
		ExpectQuery(getBooksByAuthorIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnError(self.testError)

	result, err := self.client.GetBooksByAuthorId(self.context, self.author.ID)

	self.EqualError(err, self.testError.Error())
	self.Nil(result)
}

func (self *DatabaseClientTests) TestGetBooksByAuthorIdErrorIfScanRowFailed() {
	rows := sqlmock.NewRows([]string{"id", "title", "author_id"}).
		AddRow(self.book.ID, self.book.Title, self.book.Author.ID).
		AddRow(nil, nil, nil)
	self.sqlMock.
		ExpectQuery(getBooksByAuthorIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(rows)

	result, err := self.client.GetBooksByAuthorId(self.context, self.author.ID)

	self.ErrorContains(err, "failed to scan row")
	self.Nil(result)
}

func (self *DatabaseClientTests) TestGetBooksByAuthorIdNilIfNoRows() {
	self.sqlMock.
		ExpectQuery(getBooksByAuthorIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}))

	result, err := self.client.GetBooksByAuthorId(self.context, self.author.ID)

	self.NoError(err)
	self.Nil(result)
}

func (self *DatabaseClientTests) TestGetBooksByAuthor() {
	rows := sqlmock.NewRows([]string{"id", "title", "author_id"}).
		AddRow(self.book.ID, self.book.Title, self.book.Author.ID)
	self.sqlMock.
		ExpectQuery(getBooksByAuthorIdQueryMatcher).
		WithArgs(self.author.ID).
		WillReturnRows(rows)

	result, err := self.client.GetBooksByAuthorId(self.context, self.author.ID)

	self.NoError(err)
	self.Equal([]models.Book{self.book}, result)
}
