// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	sqlx "github.com/jmoiron/sqlx"
	mock "github.com/stretchr/testify/mock"
)

// SqlxWrapper is an autogenerated mock type for the SqlxWrapper type
type SqlxWrapper struct {
	mock.Mock
}

// Connect provides a mock function with given fields: driverName, dataSourceName
func (_m *SqlxWrapper) Connect(driverName string, dataSourceName string) (*sqlx.DB, error) {
	ret := _m.Called(driverName, dataSourceName)

	var r0 *sqlx.DB
	if rf, ok := ret.Get(0).(func(string, string) *sqlx.DB); ok {
		r0 = rf(driverName, dataSourceName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.DB)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(driverName, dataSourceName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSqlxWrapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewSqlxWrapper creates a new instance of SqlxWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSqlxWrapper(t mockConstructorTestingTNewSqlxWrapper) *SqlxWrapper {
	mock := &SqlxWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
