// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateProfile mocks base method.
func (m *MockRepositoryInterface) CreateProfile(input Profile) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfile", input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfile indicates an expected call of CreateProfile.
func (mr *MockRepositoryInterfaceMockRecorder) CreateProfile(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfile", reflect.TypeOf((*MockRepositoryInterface)(nil).CreateProfile), input)
}

// GetPhoneNumberExistence mocks base method.
func (m *MockRepositoryInterface) GetPhoneNumberExistence(phoneNumber string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhoneNumberExistence", phoneNumber)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPhoneNumberExistence indicates an expected call of GetPhoneNumberExistence.
func (mr *MockRepositoryInterfaceMockRecorder) GetPhoneNumberExistence(phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhoneNumberExistence", reflect.TypeOf((*MockRepositoryInterface)(nil).GetPhoneNumberExistence), phoneNumber)
}

// GetPhoneNumberExistenceWithExcludedID mocks base method.
func (m *MockRepositoryInterface) GetPhoneNumberExistenceWithExcludedID(phoneNumber string, excludedID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhoneNumberExistenceWithExcludedID", phoneNumber, excludedID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPhoneNumberExistenceWithExcludedID indicates an expected call of GetPhoneNumberExistenceWithExcludedID.
func (mr *MockRepositoryInterfaceMockRecorder) GetPhoneNumberExistenceWithExcludedID(phoneNumber, excludedID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhoneNumberExistenceWithExcludedID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetPhoneNumberExistenceWithExcludedID), phoneNumber, excludedID)
}

// GetProfileByID mocks base method.
func (m *MockRepositoryInterface) GetProfileByID(id int) (Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByID", id)
	ret0, _ := ret[0].(Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileByID indicates an expected call of GetProfileByID.
func (mr *MockRepositoryInterfaceMockRecorder) GetProfileByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetProfileByID), id)
}

// GetProfileByPhoneNumber mocks base method.
func (m *MockRepositoryInterface) GetProfileByPhoneNumber(phoneNumber string) (Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByPhoneNumber", phoneNumber)
	ret0, _ := ret[0].(Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileByPhoneNumber indicates an expected call of GetProfileByPhoneNumber.
func (mr *MockRepositoryInterfaceMockRecorder) GetProfileByPhoneNumber(phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByPhoneNumber", reflect.TypeOf((*MockRepositoryInterface)(nil).GetProfileByPhoneNumber), phoneNumber)
}

// UpdateProfileByID mocks base method.
func (m *MockRepositoryInterface) UpdateProfileByID(profile Profile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfileByID", profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfileByID indicates an expected call of UpdateProfileByID.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateProfileByID(profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileByID", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateProfileByID), profile)
}

// UpsertProfileMetaData mocks base method.
func (m *MockRepositoryInterface) UpsertProfileMetaData(input ProfileMetaData) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProfileMetaData", input)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertProfileMetaData indicates an expected call of UpsertProfileMetaData.
func (mr *MockRepositoryInterfaceMockRecorder) UpsertProfileMetaData(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProfileMetaData", reflect.TypeOf((*MockRepositoryInterface)(nil).UpsertProfileMetaData), input)
}
