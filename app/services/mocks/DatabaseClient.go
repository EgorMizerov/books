// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "github.com/egormizerov/books/app/models"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// DatabaseClient is an autogenerated mock type for the DatabaseClient type
type DatabaseClient struct {
	mock.Mock
}

// CreateAuthor provides a mock function with given fields: ctx, author
func (_m *DatabaseClient) CreateAuthor(ctx context.Context, author models.Author) error {
	ret := _m.Called(ctx, author)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Author) error); ok {
		r0 = rf(ctx, author)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateBook provides a mock function with given fields: ctx, book
func (_m *DatabaseClient) CreateBook(ctx context.Context, book models.Book) error {
	ret := _m.Called(ctx, book)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Book) error); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAuthorById provides a mock function with given fields: ctx, authorId
func (_m *DatabaseClient) GetAuthorById(ctx context.Context, authorId uuid.UUID) (models.Author, error) {
	ret := _m.Called(ctx, authorId)

	var r0 models.Author
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Author); ok {
		r0 = rf(ctx, authorId)
	} else {
		r0 = ret.Get(0).(models.Author)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, authorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBookById provides a mock function with given fields: ctx, bookId
func (_m *DatabaseClient) GetBookById(ctx context.Context, bookId uuid.UUID) (models.Book, error) {
	ret := _m.Called(ctx, bookId)

	var r0 models.Book
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Book); ok {
		r0 = rf(ctx, bookId)
	} else {
		r0 = ret.Get(0).(models.Book)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, bookId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBooksByAuthorId provides a mock function with given fields: ctx, authorId
func (_m *DatabaseClient) GetBooksByAuthorId(ctx context.Context, authorId uuid.UUID) ([]models.Book, error) {
	ret := _m.Called(ctx, authorId)

	var r0 []models.Book
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.Book); ok {
		r0 = rf(ctx, authorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Book)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, authorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDatabaseClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewDatabaseClient creates a new instance of DatabaseClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDatabaseClient(t mockConstructorTestingTNewDatabaseClient) *DatabaseClient {
	mock := &DatabaseClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}