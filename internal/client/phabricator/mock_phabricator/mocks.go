// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -destination=mock_phabricator/mocks.go -self_package=github.com/moisesvega/diffy/internal/client/phabricator/mock_phabricator
//

// Package mock_phabricator is a generated GoMock package.
package mock_phabricator

import (
	reflect "reflect"

	phabricator "github.com/moisesvega/diffy/internal/client/phabricator"
	config "github.com/moisesvega/diffy/internal/config"
	requests "github.com/uber/gonduit/requests"
	responses "github.com/uber/gonduit/responses"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
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
func (m *MockClient) GetUsers(strings []string) ([]*phabricator.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", strings)
	ret0, _ := ret[0].([]*phabricator.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockClientMockRecorder) GetUsers(strings any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockClient)(nil).GetUsers), strings)
}

// New mocks base method.
func (m *MockClient) New(cfg *config.PhabricatorConfig) (phabricator.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", cfg)
	ret0, _ := ret[0].(phabricator.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockClientMockRecorder) New(cfg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockClient)(nil).New), cfg)
}

// Mockigonduit is a mock of igonduit interface.
type Mockigonduit struct {
	ctrl     *gomock.Controller
	recorder *MockigonduitMockRecorder
}

// MockigonduitMockRecorder is the mock recorder for Mockigonduit.
type MockigonduitMockRecorder struct {
	mock *Mockigonduit
}

// NewMockigonduit creates a new mock instance.
func NewMockigonduit(ctrl *gomock.Controller) *Mockigonduit {
	mock := &Mockigonduit{ctrl: ctrl}
	mock.recorder = &MockigonduitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockigonduit) EXPECT() *MockigonduitMockRecorder {
	return m.recorder
}

// DifferentialQuery mocks base method.
func (m *Mockigonduit) DifferentialQuery(req requests.DifferentialQueryRequest) (*responses.DifferentialQueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DifferentialQuery", req)
	ret0, _ := ret[0].(*responses.DifferentialQueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DifferentialQuery indicates an expected call of DifferentialQuery.
func (mr *MockigonduitMockRecorder) DifferentialQuery(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DifferentialQuery", reflect.TypeOf((*Mockigonduit)(nil).DifferentialQuery), req)
}

// UserQuery mocks base method.
func (m *Mockigonduit) UserQuery(req requests.UserQueryRequest) (*responses.UserQueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserQuery", req)
	ret0, _ := ret[0].(*responses.UserQueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserQuery indicates an expected call of UserQuery.
func (mr *MockigonduitMockRecorder) UserQuery(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserQuery", reflect.TypeOf((*Mockigonduit)(nil).UserQuery), req)
}