// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/Message.proto

package Message

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

// ProvClient is the client API for Prov service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProvClient interface {
	Apply(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error)
	Drop(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error)
}

type provClient struct {
	cc grpc.ClientConnInterface
}

func NewProvClient(cc grpc.ClientConnInterface) ProvClient {
	return &provClient{cc}
}

func (c *provClient) Apply(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	out := new(Resp)
	err := c.cc.Invoke(ctx, "/Prov/apply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *provClient) Drop(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	out := new(Resp)
	err := c.cc.Invoke(ctx, "/Prov/drop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProvServer is the server API for Prov service.
// All implementations must embed UnimplementedProvServer
// for forward compatibility
type ProvServer interface {
	Apply(context.Context, *Req) (*Resp, error)
	Drop(context.Context, *Req) (*Resp, error)
	mustEmbedUnimplementedProvServer()
}

// UnimplementedProvServer must be embedded to have forward compatible implementations.
type UnimplementedProvServer struct {
}

func (UnimplementedProvServer) Apply(context.Context, *Req) (*Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apply not implemented")
}
func (UnimplementedProvServer) Drop(context.Context, *Req) (*Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Drop not implemented")
}
func (UnimplementedProvServer) mustEmbedUnimplementedProvServer() {}

// UnsafeProvServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProvServer will
// result in compilation errors.
type UnsafeProvServer interface {
	mustEmbedUnimplementedProvServer()
}

func RegisterProvServer(s grpc.ServiceRegistrar, srv ProvServer) {
	s.RegisterService(&Prov_ServiceDesc, srv)
}

func _Prov_Apply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProvServer).Apply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Prov/apply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProvServer).Apply(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _Prov_Drop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProvServer).Drop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Prov/drop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProvServer).Drop(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

// Prov_ServiceDesc is the grpc.ServiceDesc for Prov service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Prov_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Prov",
	HandlerType: (*ProvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "apply",
			Handler:    _Prov_Apply_Handler,
		},
		{
			MethodName: "drop",
			Handler:    _Prov_Drop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/Message.proto",
}
