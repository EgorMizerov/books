package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/egormizerov/books/app/models"
)

var (
	ErrInvalidInputBody     = "Invalid input body."
	ErrInvalidPathVariables = "Invalid path variables."
	ErrMethodNotAllowed     = "Method not allowed."

	ErrCreateAuthor    = "We could not create new author. Please try again."
	ErrGetAuthorsBooks = "We could not get author's books. Please try again."
	ErrCreateBook      = "We could not create new book. Please try again."
	ErrGetBook         = "We could not get book. Please try again."

	EndpointCreateAuthorMatcher    = regexp.MustCompile("^/api/authors$")
	EndpointGetAuthorsBooksMatcher = regexp.MustCompile("^/api/authors/(.{36})/books/$")
	EndpointCreateBookMatcher      = regexp.MustCompile("^/api/books$")
	EndpointGetBookMatcher         = regexp.MustCompile("^/api/books/(.{36})$")
)

//go:generate mockery --name=Service
type Service interface {
	CreateAuthor(ctx context.Context, authorName string) error
	CreateBook(ctx context.Context, title string, authorId uuid.UUID) error
	GetBook(ctx context.Context, bookId uuid.UUID) (models.Book, error)
	GetAuthorsBooks(ctx context.Context, authorId uuid.UUID) ([]models.Book, error)
}

type Handler struct {
	service   Service
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewHandler(logger *logrus.Logger, service Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}

func (self *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request = MiddlewareLoggerInContext(request, self.logger)

	switch request.Method {
	case http.MethodGet:
		if EndpointGetBookMatcher.MatchString(request.URL.Path) {
			self.GetBook(response, request)
			return
		}
		if EndpointGetAuthorsBooksMatcher.MatchString(request.URL.Path) {
			self.GetAuthorsBooks(response, request)
			return
		}
	case http.MethodPost:
		if EndpointCreateAuthorMatcher.MatchString(request.URL.Path) {
			self.CreateAuthor(response, request)
			return
		}
		if EndpointCreateBookMatcher.MatchString(request.URL.Path) {
			self.CreateBook(response, request)
			return
		}
	case http.MethodOptions:
		response.Header().Set("Allow", "GET, POST, OPTIONS")
		response.WriteHeader(http.StatusNoContent)
		return
	default:
		response.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(response, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(response, request)
}

type CreateAuthorRequestBody struct {
	Name string `validate:"required"`
}

func (self *Handler) CreateAuthor(response http.ResponseWriter, request *http.Request) {
	var input CreateAuthorRequestBody
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(response, ErrCreateAuthor, http.StatusInternalServerError)
		return
	}
	if err := self.validator.Struct(input); err != nil {
		http.Error(response, ErrInvalidInputBody, http.StatusUnprocessableEntity)
		return
	}

	if err := self.service.CreateAuthor(request.Context(), input.Name); err != nil {
		http.Error(response, ErrCreateAuthor, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (self *Handler) GetAuthorsBooks(response http.ResponseWriter, request *http.Request) {
	pathComponents := EndpointGetAuthorsBooksMatcher.FindStringSubmatch(request.URL.Path)
	if len(pathComponents) < 2 || pathComponents[1] == "" {
		http.Error(response, ErrInvalidPathVariables, http.StatusUnprocessableEntity)
		return
	}
	authorId, err := uuid.Parse(pathComponents[1])
	if err != nil {
		http.Error(response, ErrGetAuthorsBooks, http.StatusInternalServerError)
		return
	}

	books, err := self.service.GetAuthorsBooks(request.Context(), authorId)
	if err != nil {
		http.Error(response, ErrGetAuthorsBooks, http.StatusInternalServerError)
		return
	}

	if books == nil {
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write([]byte("[]"))
		return
	}
	booksJson, err := json.Marshal(books)
	if err != nil {
		http.Error(response, ErrGetAuthorsBooks, http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(booksJson)
}

type CreateBookRequestBody struct {
	Title    string `json:"title" validate:"required"`
	AuthorID string `json:"author_id" validate:"required,uuid"`
}

func (self *Handler) CreateBook(response http.ResponseWriter, request *http.Request) {
	var input CreateBookRequestBody
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		http.Error(response, ErrCreateBook, http.StatusInternalServerError)
		return
	}
	if err := self.validator.Struct(input); err != nil {
		fmt.Println(input)
		http.Error(response, ErrInvalidInputBody, http.StatusUnprocessableEntity)
		return
	}

	if err := self.service.CreateBook(request.Context(), input.Title, uuid.MustParse(input.AuthorID)); err != nil {
		http.Error(response, ErrCreateBook, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (self *Handler) GetBook(response http.ResponseWriter, request *http.Request) {
	pathComponents := EndpointGetBookMatcher.FindStringSubmatch(request.URL.Path)
	if len(pathComponents) < 2 || pathComponents[1] == "" {
		http.Error(response, ErrInvalidPathVariables, http.StatusUnprocessableEntity)
		return
	}
	bookId, err := uuid.Parse(pathComponents[1])
	if err != nil {
		http.Error(response, ErrGetBook, http.StatusInternalServerError)
		return
	}

	book, err := self.service.GetBook(request.Context(), bookId)
	if err != nil {
		http.Error(response, ErrGetBook, http.StatusInternalServerError)
		return
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		http.Error(response, ErrGetBook, http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(bookJson)
}
