// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_store is a generated GoMock package.
package store

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "developer.zopsmart.com/go/gofr/examples/universal-example/cassandra/entity"
	gofr "developer.zopsmart.com/go/gofr/pkg/gofr"
)

// MockEmployee is a mock of Employee interface
type MockEmployee struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeMockRecorder
}

// MockEmployeeMockRecorder is the mock recorder for MockEmployee
type MockEmployeeMockRecorder struct {
	mock *MockEmployee
}

// NewMockEmployee creates a new mock instance
func NewMockEmployee(ctrl *gomock.Controller) *MockEmployee {
	mock := &MockEmployee{ctrl: ctrl}
	mock.recorder = &MockEmployeeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEmployee) EXPECT() *MockEmployeeMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockEmployee) Get(ctx *gofr.Context, filter entity.Employee) []entity.Employee {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, filter)
	ret0, _ := ret[0].([]entity.Employee)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockEmployeeMockRecorder) Get(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockEmployee)(nil).Get), ctx, filter)
}

// Create mocks base method
func (m *MockEmployee) Create(ctx *gofr.Context, data entity.Employee) ([]entity.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].([]entity.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockEmployeeMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmployee)(nil).Create), ctx, data)
}
