// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protoc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CountServiceClient is the client API for CountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CountServiceClient interface {
	// 客户端流式请求
	ClientData(ctx context.Context, opts ...grpc.CallOption) (CountService_ClientDataClient, error)
	// 服务端响应
	ServerCount(ctx context.Context, in *RowRequest, opts ...grpc.CallOption) (*CountResponse, error)
}

type countServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCountServiceClient(cc grpc.ClientConnInterface) CountServiceClient {
	return &countServiceClient{cc}
}

func (c *countServiceClient) ClientData(ctx context.Context, opts ...grpc.CallOption) (CountService_ClientDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &CountService_ServiceDesc.Streams[0], "/datacount.countService/ClientData", opts...)
	if err != nil {
		return nil, err
	}
	x := &countServiceClientDataClient{stream}
	return x, nil
}

type CountService_ClientDataClient interface {
	Send(*RowRequest) error
	CloseAndRecv() (*CountResponse, error)
	grpc.ClientStream
}

type countServiceClientDataClient struct {
	grpc.ClientStream
}

func (x *countServiceClientDataClient) Send(m *RowRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *countServiceClientDataClient) CloseAndRecv() (*CountResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(CountResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *countServiceClient) ServerCount(ctx context.Context, in *RowRequest, opts ...grpc.CallOption) (*CountResponse, error) {
	out := new(CountResponse)
	err := c.cc.Invoke(ctx, "/datacount.countService/ServerCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CountServiceServer is the server API for CountService service.
// All implementations must embed UnimplementedCountServiceServer
// for forward compatibility
type CountServiceServer interface {
	// 客户端流式请求
	ClientData(CountService_ClientDataServer) error
	// 服务端响应
	ServerCount(context.Context, *RowRequest) (*CountResponse, error)
	mustEmbedUnimplementedCountServiceServer()
}

// UnimplementedCountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCountServiceServer struct {
}

func (UnimplementedCountServiceServer) ClientData(CountService_ClientDataServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientData not implemented")
}
func (UnimplementedCountServiceServer) ServerCount(context.Context, *RowRequest) (*CountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServerCount not implemented")
}
func (UnimplementedCountServiceServer) mustEmbedUnimplementedCountServiceServer() {}

// UnsafeCountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CountServiceServer will
// result in compilation errors.
type UnsafeCountServiceServer interface {
	mustEmbedUnimplementedCountServiceServer()
}

func RegisterCountServiceServer(s grpc.ServiceRegistrar, srv CountServiceServer) {
	s.RegisterService(&CountService_ServiceDesc, srv)
}

func _CountService_ClientData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CountServiceServer).ClientData(&countServiceClientDataServer{stream})
}

type CountService_ClientDataServer interface {
	SendAndClose(*CountResponse) error
	Recv() (*RowRequest, error)
	grpc.ServerStream
}

type countServiceClientDataServer struct {
	grpc.ServerStream
}

func (x *countServiceClientDataServer) SendAndClose(m *CountResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *countServiceClientDataServer) Recv() (*RowRequest, error) {
	m := new(RowRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CountService_ServerCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CountServiceServer).ServerCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datacount.countService/ServerCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CountServiceServer).ServerCount(ctx, req.(*RowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CountService_ServiceDesc is the grpc.ServiceDesc for CountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datacount.countService",
	HandlerType: (*CountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServerCount",
			Handler:    _CountService_ServerCount_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ClientData",
			Handler:       _CountService_ClientData_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "datacount.proto",
}