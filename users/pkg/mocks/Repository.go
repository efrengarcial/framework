// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/efrengarcial/framework/users/pkg/model"
import paginations "github.com/efrengarcial/framework/users/pkg/utils/paginations"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, where, args
func (_m *Repository) Delete(_a0 model.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, _a0, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.IModel, string, ...interface{}) error); ok {
		r0 = rf(_a0, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: receiver, id
func (_m *Repository) Find(receiver model.IModel, id uint64) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.IModel, uint64) error); ok {
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
func (_m *Repository) FindAllPageable(pageable *model.Pageable, result interface{}, where string, args ...interface{}) (*paginations.Pagination, error) {
	var _ca []interface{}
	_ca = append(_ca, pageable, result, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *paginations.Pagination
	if rf, ok := ret.Get(0).(func(*model.Pageable, interface{}, string, ...interface{}) *paginations.Pagination); ok {
		r0 = rf(pageable, result, where, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*paginations.Pagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Pageable, interface{}, string, ...interface{}) error); ok {
		r1 = rf(pageable, result, where, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindFirst provides a mock function with given fields: receiver, where, args
func (_m *Repository) FindFirst(receiver model.IModel, where string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, receiver, where)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.IModel, string, ...interface{}) error); ok {
		r0 = rf(receiver, where, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: _a0
func (_m *Repository) Insert(_a0 model.IModel) (model.IModel, error) {
	ret := _m.Called(_a0)

	var r0 model.IModel
	if rf, ok := ret.Get(0).(func(model.IModel) model.IModel); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.IModel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRecord provides a mock function with given fields: _a0
func (_m *Repository) NewRecord(_a0 model.IModel) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(model.IModel) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *Repository) Save(_a0 model.IModel) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(model.IModel) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.IModel) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 model.IModel) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.IModel) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
