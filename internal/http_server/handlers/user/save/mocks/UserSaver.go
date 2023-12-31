// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserSaver is an autogenerated mock type for the UserSaver type
type UserSaver struct {
	mock.Mock
}

// SaveUser provides a mock function with given fields: username, email, password
func (_m *UserSaver) SaveUser(username string, email string, password string) (int64, error) {
	ret := _m.Called(username, email, password)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string) (int64, error)); ok {
		return rf(username, email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) int64); ok {
		r0 = rf(username, email, password)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(username, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserSaver creates a new instance of UserSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserSaver {
	mock := &UserSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
