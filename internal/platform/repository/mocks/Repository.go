// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, where, args
func (_m *Repository) Delete(_a0 domain.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, _a0, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.IModel, string, ...interface{}) error); ok {
		r0 = rf(_a0, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: receiver, id
func (_m *Repository) Find(receiver domain.IModel, id uint64) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.IModel, uint64) error); ok {
		r0 = rf(receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: result, where, args
func (_m *Repository) FindAll(result interface{}, where string, args ...interface{}) error {
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
func (_m *Repository) FindAllPageable(pageable *domain.Pageable, result interface{}, where string, args ...interface{}) (*database.Pagination, error) {
	var _ca []interface{}
	_ca = append(_ca, pageable, result, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *database.Pagination
	if rf, ok := ret.Get(0).(func(*domain.Pageable, interface{}, string, ...interface{}) *database.Pagination); ok {
		r0 = rf(pageable, result, where, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.Pagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.Pageable, interface{}, string, ...interface{}) error); ok {
		r1 = rf(pageable, result, where, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindFirst provides a mock function with given fields: receiver, where, args
func (_m *Repository) FindFirst(receiver domain.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, receiver, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.IModel, string, ...interface{}) error); ok {
		r0 = rf(receiver, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: _a0
func (_m *Repository) Insert(_a0 domain.IModel) (domain.IModel, error) {
	ret := _m.Called(_a0)

	var r0 domain.IModel
	if rf, ok := ret.Get(0).(func(domain.IModel) domain.IModel); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.IModel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRecord provides a mock function with given fields: _a0
func (_m *Repository) NewRecord(_a0 domain.IModel) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(domain.IModel) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *Repository) Save(_a0 domain.IModel) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(domain.IModel) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 domain.IModel) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.IModel) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
