// Code generated by MockGen. DO NOT EDIT.
// Source: ./api.go
//
// Generated by this command:
//
//	mockgen -source ./api.go -destination ./mocks/auth/auth.go -package mock_auth
//

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthenticator is a mock of Authenticator interface.
type MockAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticatorMockRecorder
	isgomock struct{}
}

// MockAuthenticatorMockRecorder is the mock recorder for MockAuthenticator.
type MockAuthenticatorMockRecorder struct {
	mock *MockAuthenticator
}

// NewMockAuthenticator creates a new mock instance.
func NewMockAuthenticator(ctrl *gomock.Controller) *MockAuthenticator {
	mock := &MockAuthenticator{ctrl: ctrl}
	mock.recorder = &MockAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticator) EXPECT() *MockAuthenticatorMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthenticator) Login(user, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", user, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthenticatorMockRecorder) Login(user, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthenticator)(nil).Login), user, password)
}
