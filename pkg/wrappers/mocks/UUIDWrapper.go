// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// UUIDWrapper is an autogenerated mock type for the UUIDWrapper type
type UUIDWrapper struct {
	mock.Mock
}

// New provides a mock function with given fields:
func (_m *UUIDWrapper) New() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

type mockConstructorTestingTNewUUIDWrapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewUUIDWrapper creates a new instance of UUIDWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUUIDWrapper(t mockConstructorTestingTNewUUIDWrapper) *UUIDWrapper {
	mock := &UUIDWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}