// Code generated by MockGen. DO NOT EDIT.
// Source: notification_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source notification_grpc.pb.go -destination mocks/mock_client_notification.go -package keeperv1 NotificationServiceClient
//

// Package keeperv1 is a generated GoMock package.
package keeperv1

import (
	context "context"
	reflect "reflect"

	keeperv1 "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockNotificationServiceClient is a mock of NotificationServiceClient interface.
type MockNotificationServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceClientMockRecorder
}

// MockNotificationServiceClientMockRecorder is the mock recorder for MockNotificationServiceClient.
type MockNotificationServiceClientMockRecorder struct {
	mock *MockNotificationServiceClient
}

// NewMockNotificationServiceClient creates a new mock instance.
func NewMockNotificationServiceClient(ctrl *gomock.Controller) *MockNotificationServiceClient {
	mock := &MockNotificationServiceClient{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceClient) EXPECT() *MockNotificationServiceClientMockRecorder {
	return m.recorder
}

// SubscribeV1 mocks base method.
func (m *MockNotificationServiceClient) SubscribeV1(ctx context.Context, in *keeperv1.SubscribeV1Request, opts ...grpc.CallOption) (grpc.ServerStreamingClient[keeperv1.SubscribeResponseV1], error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubscribeV1", varargs...)
	ret0, _ := ret[0].(grpc.ServerStreamingClient[keeperv1.SubscribeResponseV1])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeV1 indicates an expected call of SubscribeV1.
func (mr *MockNotificationServiceClientMockRecorder) SubscribeV1(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeV1", reflect.TypeOf((*MockNotificationServiceClient)(nil).SubscribeV1), varargs...)
}

// MockNotificationServiceServer is a mock of NotificationServiceServer interface.
type MockNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceServerMockRecorder
}

// MockNotificationServiceServerMockRecorder is the mock recorder for MockNotificationServiceServer.
type MockNotificationServiceServerMockRecorder struct {
	mock *MockNotificationServiceServer
}

// NewMockNotificationServiceServer creates a new mock instance.
func NewMockNotificationServiceServer(ctrl *gomock.Controller) *MockNotificationServiceServer {
	mock := &MockNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceServer) EXPECT() *MockNotificationServiceServerMockRecorder {
	return m.recorder
}

// SubscribeV1 mocks base method.
func (m *MockNotificationServiceServer) SubscribeV1(arg0 *keeperv1.SubscribeV1Request, arg1 grpc.ServerStreamingServer[keeperv1.SubscribeResponseV1]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeV1", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeV1 indicates an expected call of SubscribeV1.
func (mr *MockNotificationServiceServerMockRecorder) SubscribeV1(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeV1", reflect.TypeOf((*MockNotificationServiceServer)(nil).SubscribeV1), arg0, arg1)
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}

// MockUnsafeNotificationServiceServer is a mock of UnsafeNotificationServiceServer interface.
type MockUnsafeNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNotificationServiceServerMockRecorder
}

// MockUnsafeNotificationServiceServerMockRecorder is the mock recorder for MockUnsafeNotificationServiceServer.
type MockUnsafeNotificationServiceServerMockRecorder struct {
	mock *MockUnsafeNotificationServiceServer
}

// NewMockUnsafeNotificationServiceServer creates a new mock instance.
func NewMockUnsafeNotificationServiceServer(ctrl *gomock.Controller) *MockUnsafeNotificationServiceServer {
	mock := &MockUnsafeNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNotificationServiceServer) EXPECT() *MockUnsafeNotificationServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockUnsafeNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockUnsafeNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockUnsafeNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}
