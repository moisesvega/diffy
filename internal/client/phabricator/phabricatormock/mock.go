// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/moisesvega/diffy/internal/client/phabricator (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -destination=mock.go -package=phabricatormock -write_generate_directive github.com/moisesvega/diffy/internal/client/phabricator Client
//

// Package phabricatormock is a generated GoMock package.
package phabricatormock

import (
	reflect "reflect"

	entity "github.com/moisesvega/diffy/internal/entity"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=mock.go -package=phabricatormock -write_generate_directive github.com/moisesvega/diffy/internal/client/phabricator Client

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetUsers mocks base method.
func (m *MockClient) GetUsers(names []string) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", names)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockClientMockRecorder) GetUsers(names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockClient)(nil).GetUsers), names)
}
