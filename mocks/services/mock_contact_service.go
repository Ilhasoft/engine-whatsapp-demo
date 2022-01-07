// Code generated by MockGen. DO NOT EDIT.
// Source: services/contact_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/weni/whatsapp-router/models"
)

// MockContactService is a mock of ContactService interface.
type MockContactService struct {
	ctrl     *gomock.Controller
	recorder *MockContactServiceMockRecorder
}

// MockContactServiceMockRecorder is the mock recorder for MockContactService.
type MockContactServiceMockRecorder struct {
	mock *MockContactService
}

// NewMockContactService creates a new mock instance.
func NewMockContactService(ctrl *gomock.Controller) *MockContactService {
	mock := &MockContactService{ctrl: ctrl}
	mock.recorder = &MockContactServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactService) EXPECT() *MockContactServiceMockRecorder {
	return m.recorder
}

// CreateContact mocks base method.
func (m *MockContactService) CreateContact(arg0 *models.Contact) (*models.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContact", arg0)
	ret0, _ := ret[0].(*models.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContact indicates an expected call of CreateContact.
func (mr *MockContactServiceMockRecorder) CreateContact(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContact", reflect.TypeOf((*MockContactService)(nil).CreateContact), arg0)
}

// FindContact mocks base method.
func (m *MockContactService) FindContact(arg0 *models.Contact) (*models.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindContact", arg0)
	ret0, _ := ret[0].(*models.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindContact indicates an expected call of FindContact.
func (mr *MockContactServiceMockRecorder) FindContact(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindContact", reflect.TypeOf((*MockContactService)(nil).FindContact), arg0)
}

// UpdateContact mocks base method.
func (m *MockContactService) UpdateContact(arg0 *models.Contact) (*models.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContact", arg0)
	ret0, _ := ret[0].(*models.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateContact indicates an expected call of UpdateContact.
func (mr *MockContactServiceMockRecorder) UpdateContact(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContact", reflect.TypeOf((*MockContactService)(nil).UpdateContact), arg0)
}
