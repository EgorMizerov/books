package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"

	"github.com/egormizerov/books/app/models"
	"github.com/egormizerov/books/app/services/mocks"
	logcontext "github.com/egormizerov/books/pkg/log/context"
	wrappersmocks "github.com/egormizerov/books/pkg/wrappers/mocks"
)

type ServiceTests struct {
	suite.Suite
	service            Service
	mockDatabaseClient *mocks.DatabaseClient
	uuidMock           *wrappersmocks.UUIDWrapper
	logger             *logrus.Logger
	loggerHook         *logrustest.Hook

	testError         error
	contextWithLogger context.Context
	book              models.Book
	author            models.Author
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTests))
}

func (self *ServiceTests) SetupTest() {
	self.logger, self.loggerHook = logrustest.NewNullLogger()
	self.mockDatabaseClient = mocks.NewDatabaseClient(self.T())
	self.uuidMock = wrappersmocks.NewUUIDWrapper(self.T())
	self.service = Service{
		DatabaseClient: self.mockDatabaseClient,
		uuid:           self.uuidMock,
	}

	self.contextWithLogger = logcontext.WithLogger(context.Background(), logrus.NewEntry(self.logger))
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

func (self *ServiceTests) TestNewService() {
	result := NewService(self.mockDatabaseClient, self.uuidMock)

	self.Equal(&Service{
		DatabaseClient: self.mockDatabaseClient,
		uuid:           self.uuidMock,
	}, result)
}

func (self *ServiceTests) TestCreateBookErrorIfModelsNewBookFailed() {
	self.uuidMock.On("New").Return(self.book.ID)

	err := self.service.CreateBook(self.contextWithLogger, "", self.author.ID)

	self.ErrorContains(err, "failed to init book")
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{
			"author_id": self.author.ID.String(),
		},
		"book title must not be empty",
		"failed to init book",
	)
}

func (self *ServiceTests) TestCreateBookErrorIfCreateBookFailed() {
	self.mockDatabaseClient.
		On("CreateBook", self.contextWithLogger, self.book).
		Return(self.testError)
	self.uuidMock.On("New").Return(self.book.ID)

	err := self.service.CreateBook(self.contextWithLogger, self.book.Title, self.author.ID)

	self.ErrorContains(err, self.testError.Error())
	self.ErrorContains(err, "failed to create book")
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{
			"author_id":  self.author.ID.String(),
			"book_title": self.book.Title,
		},
		self.testError.Error(),
		"failed to create book",
	)
}

func (self *ServiceTests) TestCreateBook() {
	self.mockDatabaseClient.
		On("CreateBook", self.contextWithLogger, self.book).
		Return(nil)
	self.uuidMock.On("New").Return(self.book.ID)

	err := self.service.CreateBook(self.contextWithLogger, self.book.Title, self.author.ID)

	self.NoError(err)
}

func (self *ServiceTests) TestCreateAuthorErrorIfModelsNewAuthorFailed() {
	self.uuidMock.On("New").Return(self.author.ID)

	err := self.service.CreateAuthor(self.contextWithLogger, "")

	self.ErrorContains(err, "failed to init author")
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{},
		"author name must not be empty",
		"failed to init author",
	)
}

func (self *ServiceTests) TestCreateAuthorErrorIfCreateAuthorFailed() {
	self.mockDatabaseClient.
		On("CreateAuthor", self.contextWithLogger, self.author).
		Return(self.testError)
	self.uuidMock.On("New").Return(self.author.ID)

	err := self.service.CreateAuthor(self.contextWithLogger, self.author.Name)

	self.ErrorContains(err, self.testError.Error())
	self.ErrorContains(err, "failed to create author")
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{
			"author_name": self.author.Name,
		},
		self.testError.Error(),
		"failed to create author",
	)
}

func (self *ServiceTests) TestCreateAuthor() {
	self.mockDatabaseClient.
		On("CreateAuthor", self.contextWithLogger, self.author).
		Return(nil)
	self.uuidMock.On("New").Return(self.author.ID)

	err := self.service.CreateAuthor(self.contextWithLogger, self.author.Name)

	self.NoError(err)
}

func (self *ServiceTests) TestGetBookErrorIfGetBookByIdFailed() {
	self.mockDatabaseClient.
		On("GetBookById", self.contextWithLogger, self.book.ID).
		Return(models.Book{}, self.testError)

	result, err := self.service.GetBook(self.contextWithLogger, self.book.ID)

	self.ErrorContains(err, "failed to get book by id")
	self.Equal(models.Book{}, result)
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{
			"book_id": self.book.ID.String(),
		},
		self.testError.Error(),
		"failed to get book",
	)
}

func (self *ServiceTests) TestGetBookErrorIfGetAuthorByIdFailed() {
	self.mockDatabaseClient.
		On("GetBookById", self.contextWithLogger, self.book.ID).
		Return(self.book, nil)
	self.mockDatabaseClient.
		On("GetAuthorById", self.contextWithLogger, self.author.ID).
		Return(models.Author{}, self.testError)

	result, err := self.service.GetBook(self.contextWithLogger, self.book.ID)

	self.ErrorContains(err, self.testError.Error())
	self.ErrorContains(err, "failed to get author by id")
	self.Equal(models.Book{}, result)
	self.matchLogWithError(
		self.loggerHook.LastEntry(),
		logrus.Fields{
			"book_id": self.book.ID.String(),
		},
		self.testError.Error(),
		"failed to get author",
	)
}

func (self *ServiceTests) TestGetBook() {
	self.mockDatabaseClient.
		On("GetBookById", self.contextWithLogger, self.book.ID).
		Return(self.book, nil)
	self.mockDatabaseClient.
		On("GetAuthorById", self.contextWithLogger, self.author.ID).
		Return(self.author, nil)

	result, err := self.service.GetBook(self.contextWithLogger, self.book.ID)

	self.NoError(err)
	self.Equal(models.Book{
		ID:     self.book.ID,
		Title:  self.book.Title,
		Author: self.author,
	}, result)
}

func (self *ServiceTests) TestGetAuthorsBooksErrorIfGetBooksByAuthorIdFailed() {
	self.mockDatabaseClient.
		On("GetBooksByAuthorId", self.contextWithLogger, self.author.ID).
		Return(nil, self.testError)

	result, err := self.service.GetAuthorsBooks(self.contextWithLogger, self.author.ID)

	self.ErrorContains(err, self.testError.Error())
	self.ErrorContains(err, "failed to get books by author id")
	self.Nil(result)
}

func (self *ServiceTests) TestGetAuthorsBooks() {
	self.mockDatabaseClient.
		On("GetBooksByAuthorId", self.contextWithLogger, self.author.ID).
		Return([]models.Book{self.book}, nil)

	result, err := self.service.GetAuthorsBooks(self.contextWithLogger, self.author.ID)

	self.NoError(err)
	self.Equal([]models.Book{self.book}, result)
}

func (self *ServiceTests) matchLogWithError(
	entry *logrus.Entry,
	fields logrus.Fields,
	errorMessage string,
	logMessage string,
) {
	self.Equal(logrus.ErrorLevel, entry.Level)
	self.Contains(entry.Message, logMessage)
	self.Equal(len(fields)+1, len(entry.Data))
	self.Contains(entry.Data[logrus.ErrorKey].(error).Error(), errorMessage)

	for field, value := range fields {
		self.Equal(value, entry.Data[field])
	}
}
