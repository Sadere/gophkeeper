// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/keeper/v1/secrets.proto

package keeperv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SecretsService_SecretPreviewsV1_FullMethodName = "/proto.keeper.v1.SecretsService/SecretPreviewsV1"
	SecretsService_SaveUserSecretV1_FullMethodName = "/proto.keeper.v1.SecretsService/SaveUserSecretV1"
)

// SecretsServiceClient is the client API for SecretsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretsServiceClient interface {
	SecretPreviewsV1(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SecretPreviewsV1Response, error)
	SaveUserSecretV1(ctx context.Context, in *SaveUserSecretV1Request, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type secretsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretsServiceClient(cc grpc.ClientConnInterface) SecretsServiceClient {
	return &secretsServiceClient{cc}
}

func (c *secretsServiceClient) SecretPreviewsV1(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SecretPreviewsV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SecretPreviewsV1Response)
	err := c.cc.Invoke(ctx, SecretsService_SecretPreviewsV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretsServiceClient) SaveUserSecretV1(ctx context.Context, in *SaveUserSecretV1Request, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SecretsService_SaveUserSecretV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretsServiceServer is the server API for SecretsService service.
// All implementations must embed UnimplementedSecretsServiceServer
// for forward compatibility.
type SecretsServiceServer interface {
	SecretPreviewsV1(context.Context, *emptypb.Empty) (*SecretPreviewsV1Response, error)
	SaveUserSecretV1(context.Context, *SaveUserSecretV1Request) (*emptypb.Empty, error)
	mustEmbedUnimplementedSecretsServiceServer()
}

// UnimplementedSecretsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSecretsServiceServer struct{}

func (UnimplementedSecretsServiceServer) SecretPreviewsV1(context.Context, *emptypb.Empty) (*SecretPreviewsV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretPreviewsV1 not implemented")
}
func (UnimplementedSecretsServiceServer) SaveUserSecretV1(context.Context, *SaveUserSecretV1Request) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveUserSecretV1 not implemented")
}
func (UnimplementedSecretsServiceServer) mustEmbedUnimplementedSecretsServiceServer() {}
func (UnimplementedSecretsServiceServer) testEmbeddedByValue()                        {}

// UnsafeSecretsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretsServiceServer will
// result in compilation errors.
type UnsafeSecretsServiceServer interface {
	mustEmbedUnimplementedSecretsServiceServer()
}

func RegisterSecretsServiceServer(s grpc.ServiceRegistrar, srv SecretsServiceServer) {
	// If the following call pancis, it indicates UnimplementedSecretsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SecretsService_ServiceDesc, srv)
}

func _SecretsService_SecretPreviewsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServiceServer).SecretPreviewsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SecretsService_SecretPreviewsV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServiceServer).SecretPreviewsV1(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretsService_SaveUserSecretV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveUserSecretV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServiceServer).SaveUserSecretV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SecretsService_SaveUserSecretV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServiceServer).SaveUserSecretV1(ctx, req.(*SaveUserSecretV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretsService_ServiceDesc is the grpc.ServiceDesc for SecretsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.keeper.v1.SecretsService",
	HandlerType: (*SecretsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SecretPreviewsV1",
			Handler:    _SecretsService_SecretPreviewsV1_Handler,
		},
		{
			MethodName: "SaveUserSecretV1",
			Handler:    _SecretsService_SaveUserSecretV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/keeper/v1/secrets.proto",
}