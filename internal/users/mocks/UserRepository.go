// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/users"
	service2 "github.com/efrengarcial/framework/internal/users/service"
	"github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, where, args
func (_m *UserRepository) Delete(_a0 service2.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, _a0, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(service2.IModel, string, ...interface{}) error); ok {
		r0 = rf(_a0, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: receiver, id
func (_m *UserRepository) Find(receiver service2.IModel, id uint64) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(service2.IModel, uint64) error); ok {
		r0 = rf(receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: result, where, args
func (_m *UserRepository) FindAll(result interface{}, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, result, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...interface{}) error); ok {
		r0 = rf(result, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllPageable provides a mock function with given fields: pageable, result, where, args
func (_m *UserRepository) FindAllPageable(pageable *service2.Pageable, result interface{}, where string, args ...interface{}) (*database.Pagination, error) {
	var _ca []interface{}
	_ca = append(_ca, pageable, result, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *database.Pagination
	if rf, ok := ret.Get(0).(func(*service2.Pageable, interface{}, string, ...interface{}) *database.Pagination); ok {
		r0 = rf(pageable, result, where, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.Pagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*service2.Pageable, interface{}, string, ...interface{}) error); ok {
		r1 = rf(pageable, result, where, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindFirst provides a mock function with given fields: receiver, where, args
func (_m *UserRepository) FindFirst(receiver service2.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, receiver, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(service2.IModel, string, ...interface{}) error); ok {
		r0 = rf(receiver, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneByEmail provides a mock function with given fields: login
func (_m *UserRepository) FindOneByEmail(login string) (*users.User, error) {
	ret := _m.Called(login)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(string) *users.User); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByLogin provides a mock function with given fields: login
func (_m *UserRepository) FindOneByLogin(login string) (*users.User, error) {
	ret := _m.Called(login)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(string) *users.User); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: email
func (_m *UserRepository) GetByEmail(email string) (*users.User, error) {
	ret := _m.Called(email)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(string) *users.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByLogin provides a mock function with given fields: login
func (_m *UserRepository) GetByLogin(login string) (*users.User, error) {
	ret := _m.Called(login)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(string) *users.User); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0
func (_m *UserRepository) Insert(_a0 service2.IModel) (service2.IModel, error) {
	ret := _m.Called(_a0)

	var r0 service2.IModel
	if rf, ok := ret.Get(0).(func(service2.IModel) service2.IModel); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service2.IModel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(service2.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRecord provides a mock function with given fields: _a0
func (_m *UserRepository) NewRecord(_a0 service2.IModel) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(service2.IModel) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *UserRepository) Save(_a0 service2.IModel) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(service2.IModel) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(service2.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *UserRepository) Update(_a0 service2.IModel) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(service2.IModel) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
