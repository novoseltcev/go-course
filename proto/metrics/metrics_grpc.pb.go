// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: metrics.proto

package metrics

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
	MetricsService_GetOne_FullMethodName      = "/metrics.MetricsService/GetOne"
	MetricsService_GetAll_FullMethodName      = "/metrics.MetricsService/GetAll"
	MetricsService_Update_FullMethodName      = "/metrics.MetricsService/Update"
	MetricsService_UpdateBatch_FullMethodName = "/metrics.MetricsService/UpdateBatch"
)

// MetricsServiceClient is the client API for MetricsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsServiceClient interface {
	GetOne(ctx context.Context, in *GetOneRequest, opts ...grpc.CallOption) (*GetOneResponse, error)
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	UpdateBatch(ctx context.Context, in *UpdateBatchRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error)
}

type metricsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsServiceClient(cc grpc.ClientConnInterface) MetricsServiceClient {
	return &metricsServiceClient{cc}
}

func (c *metricsServiceClient) GetOne(ctx context.Context, in *GetOneRequest, opts ...grpc.CallOption) (*GetOneResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOneResponse)
	err := c.cc.Invoke(ctx, MetricsService_GetOne_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, MetricsService_GetAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, MetricsService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) UpdateBatch(ctx context.Context, in *UpdateBatchRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateBatchResponse)
	err := c.cc.Invoke(ctx, MetricsService_UpdateBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServiceServer is the server API for MetricsService service.
// All implementations must embed UnimplementedMetricsServiceServer
// for forward compatibility.
type MetricsServiceServer interface {
	GetOne(context.Context, *GetOneRequest) (*GetOneResponse, error)
	GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	UpdateBatch(context.Context, *UpdateBatchRequest) (*UpdateBatchResponse, error)
	mustEmbedUnimplementedMetricsServiceServer()
}

// UnimplementedMetricsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServiceServer struct{}

func (UnimplementedMetricsServiceServer) GetOne(context.Context, *GetOneRequest) (*GetOneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedMetricsServiceServer) GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedMetricsServiceServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMetricsServiceServer) UpdateBatch(context.Context, *UpdateBatchRequest) (*UpdateBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBatch not implemented")
}
func (UnimplementedMetricsServiceServer) mustEmbedUnimplementedMetricsServiceServer() {}
func (UnimplementedMetricsServiceServer) testEmbeddedByValue()                        {}

// UnsafeMetricsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServiceServer will
// result in compilation errors.
type UnsafeMetricsServiceServer interface {
	mustEmbedUnimplementedMetricsServiceServer()
}

func RegisterMetricsServiceServer(s grpc.ServiceRegistrar, srv MetricsServiceServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MetricsService_ServiceDesc, srv)
}

func _MetricsService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_GetOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetOne(ctx, req.(*GetOneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_GetAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetAll(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_UpdateBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).UpdateBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_UpdateBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).UpdateBatch(ctx, req.(*UpdateBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsService_ServiceDesc is the grpc.ServiceDesc for MetricsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _MetricsService_GetOne_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _MetricsService_GetAll_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _MetricsService_Update_Handler,
		},
		{
			MethodName: "UpdateBatch",
			Handler:    _MetricsService_UpdateBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}
