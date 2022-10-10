package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/egormizerov/books/app/handlers/mocks"
	"github.com/egormizerov/books/app/models"
	logcontext "github.com/egormizerov/books/pkg/log/context"
)

var (
	EndpointCreateAuthor    = "/api/authors"
	EndpointGetAuthorsBooks = "/api/authors/%s/books/"
	EndpointCreateBook      = "/api/books"
	EndpointGetBook         = "/api/books/%s"
)

type HandlerTests struct {
	suite.Suite
	handler     Handler
	serviceMock *mocks.Service
	validator   *validator.Validate
	logger      *logrus.Logger
	book        models.Book
	author      models.Author
	testError   error
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTests))
}

func (self *HandlerTests) SetupTest() {
	self.serviceMock = mocks.NewService(self.T())
	self.logger = logrus.New()
	self.validator = validator.New()
	self.handler = Handler{
		service:   self.serviceMock,
		logger:    self.logger,
		validator: self.validator,
	}
	self.author = models.Author{
		ID:   uuid.New(),
		Name: "test_name",
	}
	self.book = models.Book{
		ID:     uuid.New(),
		Title:  "test_title",
		Author: self.author,
	}
	self.testError = errors.New("test_error")
}

func (self *HandlerTests) TestNewHandler() {
	result := NewHandler(self.logger, self.serviceMock, self.validator)

	self.Equal(&Handler{
		service:   self.serviceMock,
		logger:    self.logger,
		validator: self.validator,
	}, result)
}

func (self *HandlerTests) TestServeHTTPCreateBook() {
	requestBody := CreateBookRequestBody{
		Title:    self.book.Title,
		AuthorID: self.book.Author.ID.String(),
	}
	response, request := self.getRequestAndResponse(http.MethodPost, EndpointCreateBook, requestBody)
	self.serviceMock.
		On("CreateBook", self.requestWithLogger(request).Context(), requestBody.Title, uuid.MustParse(requestBody.AuthorID)).
		Return(nil)

	self.handler.ServeHTTP(response, request)

	self.Equal(http.StatusCreated, response.Code)
}

func (self *HandlerTests) TestServeHTTPCreateAuthor() {
	requestBody := CreateAuthorRequestBody{
		Name: self.author.Name,
	}
	response, request := self.getRequestAndResponse(http.MethodPost, EndpointCreateAuthor, requestBody)
	self.serviceMock.
		On("CreateAuthor", self.requestWithLogger(request).Context(), requestBody.Name).
		Return(nil)

	self.handler.ServeHTTP(response, request)

	self.Equal(http.StatusCreated, response.Code)
}

func (self *HandlerTests) TestServeHTTPGetBook() {
	requestEndpoint := fmt.Sprintf(EndpointGetBook, self.book.ID.String())
	response, request := self.getRequestAndResponse(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetBook", self.requestWithLogger(request).Context(), self.book.ID).
		Return(self.book, nil)

	self.handler.ServeHTTP(response, request)

	self.Equal(http.StatusOK, response.Code)
	self.Contains(response.Body.String(), string(self.mustMarshal(self.book)))
}

func (self *HandlerTests) TestServeHTTPGetAuthorsBooks() {
	books := []models.Book{self.book}
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, self.author.ID.String())
	response, request := self.getRequestAndResponse(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetAuthorsBooks", self.requestWithLogger(request).Context(), self.author.ID).
		Return(books, nil)

	self.handler.ServeHTTP(response, request)

	self.Equal(http.StatusOK, response.Code)
	self.Contains(response.Body.String(), string(self.mustMarshal(books)))
}

func (self *HandlerTests) TestServerHTTPOptions() {
	response, request := self.getRequestAndResponse(http.MethodOptions, "/", nil)

	self.handler.ServeHTTP(response, request)

	self.Equal("GET, POST, OPTIONS", response.Header().Get("Allow"))
	self.Equal(http.StatusNoContent, response.Code)
}

func (self *HandlerTests) TestServerHTTPErrorIfMethodNotAllowed() {
	response, request := self.getRequestAndResponse(http.MethodPut, "/", nil)

	self.handler.ServeHTTP(response, request)

	self.Equal("GET, POST, OPTIONS", response.Header().Get("Allow"))
	self.Equal(http.StatusMethodNotAllowed, response.Code)
	self.Contains(response.Body.String(), ErrMethodNotAllowed)
}

func (self *HandlerTests) TestServerHTTPNotFound() {
	response, request := self.getRequestAndResponse(http.MethodGet, "/", nil)

	self.handler.ServeHTTP(response, request)

	self.Equal(http.StatusNotFound, response.Code)
	self.Contains(response.Body.String(), "404 page not found")
}

func (self *HandlerTests) TestCreateAuthorErrorIfJsonDecodeFailed() {
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateAuthor, "")

	self.handler.CreateAuthor(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrCreateAuthor)
}

func (self *HandlerTests) TestCreateAuthorErrorIfValidateFailed() {
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateAuthor, nil)

	self.handler.CreateAuthor(response, request)

	self.Equal(http.StatusUnprocessableEntity, response.Code)
	self.Contains(response.Body.String(), ErrInvalidInputBody)
}

func (self *HandlerTests) TestCreateAuthorErrorIfServiceFailed() {
	requestBody := CreateAuthorRequestBody{
		Name: "test_name",
	}
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateAuthor, requestBody)
	self.serviceMock.
		On("CreateAuthor", request.Context(), requestBody.Name).
		Return(self.testError)

	self.handler.CreateAuthor(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrCreateAuthor)
}

func (self *HandlerTests) TestCreateAuthor() {
	requestBody := CreateAuthorRequestBody{
		Name: "test_name",
	}
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateAuthor, requestBody)
	self.serviceMock.
		On("CreateAuthor", request.Context(), requestBody.Name).
		Return(nil)

	self.handler.CreateAuthor(response, request)

	self.Equal(http.StatusCreated, response.Code)
}

func (self *HandlerTests) TestGetAuthorsBooksErrorIfInvalidInput() {
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, "")
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)

	self.handler.GetAuthorsBooks(response, request)

	self.Equal(http.StatusUnprocessableEntity, response.Code)
	self.Contains(response.Body.String(), ErrInvalidPathVariables)
}

func (self *HandlerTests) TestGetAuthorsBooksErrorIfParseUUIDFailed() {
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, "b9bac125+94e7+4c4b+8df2+6cd055402bcc")
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)

	self.handler.GetAuthorsBooks(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrGetAuthorsBooks)
}

func (self *HandlerTests) TestGetAuthorsBooksErrorIfGetAuthorsBooksFailed() {
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, self.author.ID.String())
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetAuthorsBooks", request.Context(), self.author.ID).
		Return(nil, self.testError)

	self.handler.GetAuthorsBooks(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrGetAuthorsBooks)
}

func (self *HandlerTests) TestGetAuthorsBooksIfGetAuthorsBooksReturnsEmptyBooks() {
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, self.author.ID.String())
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetAuthorsBooks", request.Context(), self.author.ID).
		Return([]models.Book{}, nil)

	self.handler.GetAuthorsBooks(response, request)

	self.Equal(http.StatusOK, response.Code)
	self.Contains(response.Body.String(), "[]")
}

func (self *HandlerTests) TestGetAuthorsBooks() {
	books := []models.Book{self.book}
	requestEndpoint := fmt.Sprintf(EndpointGetAuthorsBooks, self.author.ID.String())
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetAuthorsBooks", request.Context(), self.author.ID).
		Return(books, nil)

	self.handler.GetAuthorsBooks(response, request)

	self.Equal(http.StatusOK, response.Code)
	self.Contains(response.Body.String(), string(self.mustMarshal(books)))
}

func (self *HandlerTests) TestCreateBookErrorIfJsonDecodeFailed() {
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateBook, "")

	self.handler.CreateBook(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrCreateBook)
}

func (self *HandlerTests) TestCreateBookErrorIfValidateFailed() {
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateBook, nil)

	self.handler.CreateBook(response, request)

	self.Equal(http.StatusUnprocessableEntity, response.Code)
	self.Contains(response.Body.String(), ErrInvalidInputBody)
}

func (self *HandlerTests) TestCreateBookErrorIfServiceFailed() {
	requestBody := CreateBookRequestBody{
		Title:    self.book.Title,
		AuthorID: self.book.Author.ID.String(),
	}
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateBook, requestBody)
	self.serviceMock.
		On("CreateBook", request.Context(), requestBody.Title, uuid.MustParse(requestBody.AuthorID)).
		Return(self.testError)

	self.handler.CreateBook(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrCreateBook)
}

func (self *HandlerTests) TestCreateBook() {
	requestBody := CreateBookRequestBody{
		Title:    self.book.Title,
		AuthorID: self.book.Author.ID.String(),
	}
	response, request := self.getRequestAndResponseWithLogger(http.MethodPost, EndpointCreateBook, requestBody)
	self.serviceMock.
		On("CreateBook", request.Context(), requestBody.Title, uuid.MustParse(requestBody.AuthorID)).
		Return(nil)

	self.handler.CreateBook(response, request)

	self.Equal(http.StatusCreated, response.Code)
}

func (self *HandlerTests) TestGetBookErrorIfInvalidInput() {
	requestEndpoint := fmt.Sprintf(EndpointGetBook, "")
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)

	self.handler.GetBook(response, request)

	self.Equal(http.StatusUnprocessableEntity, response.Code)
	self.Contains(response.Body.String(), ErrInvalidPathVariables)
}

func (self *HandlerTests) TestGetBookErrorIfParseUUIDFailed() {
	requestEndpoint := fmt.Sprintf(EndpointGetBook, "b9bac125+94e7+4c4b+8df2+6cd055402bcc")
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)

	self.handler.GetBook(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrGetBook)
}

func (self *HandlerTests) TestGetBookErrorIfServiceFailed() {
	requestEndpoint := fmt.Sprintf(EndpointGetBook, self.book.ID.String())
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetBook", request.Context(), self.book.ID).
		Return(models.Book{}, self.testError)

	self.handler.GetBook(response, request)

	self.Equal(http.StatusInternalServerError, response.Code)
	self.Contains(response.Body.String(), ErrGetBook)
}

func (self *HandlerTests) TestGetBook() {
	requestEndpoint := fmt.Sprintf(EndpointGetBook, self.book.ID.String())
	response, request := self.getRequestAndResponseWithLogger(http.MethodGet, requestEndpoint, nil)
	self.serviceMock.
		On("GetBook", request.Context(), self.book.ID).
		Return(self.book, nil)

	self.handler.GetBook(response, request)

	self.Equal(http.StatusOK, response.Code)
	self.Contains(response.Body.String(), string(self.mustMarshal(self.book)))
}

func (self *HandlerTests) getRequestAndResponse(httpMethod string, endpoint string, body any) (*httptest.ResponseRecorder, *http.Request) {
	requestBodyReader := bytes.NewReader(self.mustMarshal(body))
	request := httptest.NewRequest(httpMethod, endpoint, requestBodyReader)
	response := httptest.NewRecorder()
	return response, request
}

func (self *HandlerTests) getRequestAndResponseWithLogger(httpMethod string, endpoint string, body any) (*httptest.ResponseRecorder, *http.Request) {
	response, request := self.getRequestAndResponse(httpMethod, endpoint, body)
	request = self.requestWithLogger(request)
	return response, request
}

func (self *HandlerTests) requestWithLogger(request *http.Request) *http.Request {
	return request.WithContext(
		logcontext.WithLogger(
			request.Context(),
			logrus.NewEntry(self.logger),
		),
	)
}

func (self *HandlerTests) mustMarshal(value any) []byte {
	data, err := json.Marshal(value)
	self.NoError(err)
	return data
}
