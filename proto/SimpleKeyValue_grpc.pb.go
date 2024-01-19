// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: proto/SimpleKeyValue.proto

package SimpleKeyValue

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

// SimpleKeyValueClient is the client API for SimpleKeyValue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleKeyValueClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type simpleKeyValueClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleKeyValueClient(cc grpc.ClientConnInterface) SimpleKeyValueClient {
	return &simpleKeyValueClient{cc}
}

func (c *simpleKeyValueClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, "/simplekv.SimpleKeyValue/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simpleKeyValueClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/simplekv.SimpleKeyValue/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleKeyValueServer is the server API for SimpleKeyValue service.
// All implementations must embed UnimplementedSimpleKeyValueServer
// for forward compatibility
type SimpleKeyValueServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedSimpleKeyValueServer()
}

// UnimplementedSimpleKeyValueServer must be embedded to have forward compatible implementations.
type UnimplementedSimpleKeyValueServer struct {
}

func (UnimplementedSimpleKeyValueServer) Put(context.Context, *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedSimpleKeyValueServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSimpleKeyValueServer) mustEmbedUnimplementedSimpleKeyValueServer() {}

// UnsafeSimpleKeyValueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleKeyValueServer will
// result in compilation errors.
type UnsafeSimpleKeyValueServer interface {
	mustEmbedUnimplementedSimpleKeyValueServer()
}

func RegisterSimpleKeyValueServer(s grpc.ServiceRegistrar, srv SimpleKeyValueServer) {
	s.RegisterService(&SimpleKeyValue_ServiceDesc, srv)
}

func _SimpleKeyValue_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleKeyValueServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simplekv.SimpleKeyValue/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleKeyValueServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimpleKeyValue_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleKeyValueServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/simplekv.SimpleKeyValue/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleKeyValueServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimpleKeyValue_ServiceDesc is the grpc.ServiceDesc for SimpleKeyValue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimpleKeyValue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "simplekv.SimpleKeyValue",
	HandlerType: (*SimpleKeyValueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _SimpleKeyValue_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SimpleKeyValue_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/SimpleKeyValue.proto",
}
