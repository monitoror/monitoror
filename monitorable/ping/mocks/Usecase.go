// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bind "github.com/jsdidierlaurent/monitoror/pkg/bind"
import mock "github.com/stretchr/testify/mock"

import tiles "github.com/jsdidierlaurent/monitoror/models/tiles"

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Ping provides a mock function with given fields: binder
func (_m *Usecase) Ping(binder bind.Binder) (*tiles.HealthTile, error) {
	ret := _m.Called(binder)

	var r0 *tiles.HealthTile
	if rf, ok := ret.Get(0).(func(bind.Binder) *tiles.HealthTile); ok {
		r0 = rf(binder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tiles.HealthTile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bind.Binder) error); ok {
		r1 = rf(binder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
