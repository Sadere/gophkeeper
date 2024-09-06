// Code generated by MockGen. DO NOT EDIT.
// Source: secret.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/Sadere/gophkeeper/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MockSecretRepository is a mock of SecretRepository interface.
type MockSecretRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSecretRepositoryMockRecorder
}

// MockSecretRepositoryMockRecorder is the mock recorder for MockSecretRepository.
type MockSecretRepositoryMockRecorder struct {
	mock *MockSecretRepository
}

// NewMockSecretRepository creates a new mock instance.
func NewMockSecretRepository(ctrl *gomock.Controller) *MockSecretRepository {
	mock := &MockSecretRepository{ctrl: ctrl}
	mock.recorder = &MockSecretRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretRepository) EXPECT() *MockSecretRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSecretRepository) Create(ctx context.Context, secret *model.Secret) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, secret)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSecretRepositoryMockRecorder) Create(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSecretRepository)(nil).Create), ctx, secret)
}

// GetSecret mocks base method.
func (m *MockSecretRepository) GetSecret(ctx context.Context, secretID, userID uint64) (*model.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, secretID, userID)
	ret0, _ := ret[0].(*model.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretRepositoryMockRecorder) GetSecret(ctx, secretID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretRepository)(nil).GetSecret), ctx, secretID, userID)
}

// GetUserSecrets mocks base method.
func (m *MockSecretRepository) GetUserSecrets(ctx context.Context, userID uint64) (model.Secrets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSecrets", ctx, userID)
	ret0, _ := ret[0].(model.Secrets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSecrets indicates an expected call of GetUserSecrets.
func (mr *MockSecretRepositoryMockRecorder) GetUserSecrets(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSecrets", reflect.TypeOf((*MockSecretRepository)(nil).GetUserSecrets), ctx, userID)
}

// Update mocks base method.
func (m *MockSecretRepository) Update(ctx context.Context, secret *model.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, secret)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSecretRepositoryMockRecorder) Update(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSecretRepository)(nil).Update), ctx, secret)
}
