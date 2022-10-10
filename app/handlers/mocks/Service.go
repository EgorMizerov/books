// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/egormizerov/books/app/models"

	uuid "github.com/google/uuid"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateAuthor provides a mock function with given fields: ctx, authorName
func (_m *Service) CreateAuthor(ctx context.Context, authorName string) error {
	ret := _m.Called(ctx, authorName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, authorName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateBook provides a mock function with given fields: ctx, title, authorId
func (_m *Service) CreateBook(ctx context.Context, title string, authorId uuid.UUID) error {
	ret := _m.Called(ctx, title, authorId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) error); ok {
		r0 = rf(ctx, title, authorId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAuthorsBooks provides a mock function with given fields: ctx, authorId
func (_m *Service) GetAuthorsBooks(ctx context.Context, authorId uuid.UUID) ([]models.Book, error) {
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

// GetBook provides a mock function with given fields: ctx, bookId
func (_m *Service) GetBook(ctx context.Context, bookId uuid.UUID) (models.Book, error) {
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

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
