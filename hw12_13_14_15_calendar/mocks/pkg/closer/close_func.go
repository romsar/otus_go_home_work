// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CloseFunc is an autogenerated mock type for the CloseFunc type
type CloseFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *CloseFunc) Execute() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCloseFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewCloseFunc creates a new instance of CloseFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCloseFunc(t mockConstructorTestingTNewCloseFunc) *CloseFunc {
	mock := &CloseFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}