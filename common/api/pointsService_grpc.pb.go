// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: api/pointsService.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PointsService_GetUserPointsPage_FullMethodName = "/api.PointsService/GetUserPointsPage"
)

// PointsServiceClient is the client API for PointsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PointsServiceClient interface {
	GetUserPointsPage(ctx context.Context, in *PointsPageRequest, opts ...grpc.CallOption) (*PointsPage, error)
}

type pointsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPointsServiceClient(cc grpc.ClientConnInterface) PointsServiceClient {
	return &pointsServiceClient{cc}
}

func (c *pointsServiceClient) GetUserPointsPage(ctx context.Context, in *PointsPageRequest, opts ...grpc.CallOption) (*PointsPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsPage)
	err := c.cc.Invoke(ctx, PointsService_GetUserPointsPage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PointsServiceServer is the server API for PointsService service.
// All implementations must embed UnimplementedPointsServiceServer
// for forward compatibility.
type PointsServiceServer interface {
	GetUserPointsPage(context.Context, *PointsPageRequest) (*PointsPage, error)
	mustEmbedUnimplementedPointsServiceServer()
}

// UnimplementedPointsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPointsServiceServer struct{}

func (UnimplementedPointsServiceServer) GetUserPointsPage(context.Context, *PointsPageRequest) (*PointsPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPointsPage not implemented")
}
func (UnimplementedPointsServiceServer) mustEmbedUnimplementedPointsServiceServer() {}
func (UnimplementedPointsServiceServer) testEmbeddedByValue()                       {}

// UnsafePointsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PointsServiceServer will
// result in compilation errors.
type UnsafePointsServiceServer interface {
	mustEmbedUnimplementedPointsServiceServer()
}

func RegisterPointsServiceServer(s grpc.ServiceRegistrar, srv PointsServiceServer) {
	// If the following call pancis, it indicates UnimplementedPointsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PointsService_ServiceDesc, srv)
}

func _PointsService_GetUserPointsPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PointsPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PointsServiceServer).GetUserPointsPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PointsService_GetUserPointsPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PointsServiceServer).GetUserPointsPage(ctx, req.(*PointsPageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PointsService_ServiceDesc is the grpc.ServiceDesc for PointsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PointsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.PointsService",
	HandlerType: (*PointsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserPointsPage",
			Handler:    _PointsService_GetUserPointsPage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/pointsService.proto",
}