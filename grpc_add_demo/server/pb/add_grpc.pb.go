// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: pb/add.proto

package pb

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

// CalcuClient is the client API for Calcu service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalcuClient interface {
	Add(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error)
}

type calcuClient struct {
	cc grpc.ClientConnInterface
}

func NewCalcuClient(cc grpc.ClientConnInterface) CalcuClient {
	return &calcuClient{cc}
}

func (c *calcuClient) Add(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error) {
	out := new(Output)
	err := c.cc.Invoke(ctx, "/pb.Calcu/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalcuServer is the server API for Calcu service.
// All implementations must embed UnimplementedCalcuServer
// for forward compatibility
type CalcuServer interface {
	Add(context.Context, *Input) (*Output, error)
	mustEmbedUnimplementedCalcuServer()
}

// UnimplementedCalcuServer must be embedded to have forward compatible implementations.
type UnimplementedCalcuServer struct {
}

func (UnimplementedCalcuServer) Add(context.Context, *Input) (*Output, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedCalcuServer) mustEmbedUnimplementedCalcuServer() {}

// UnsafeCalcuServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalcuServer will
// result in compilation errors.
type UnsafeCalcuServer interface {
	mustEmbedUnimplementedCalcuServer()
}

func RegisterCalcuServer(s grpc.ServiceRegistrar, srv CalcuServer) {
	s.RegisterService(&Calcu_ServiceDesc, srv)
}

func _Calcu_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcuServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Calcu/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcuServer).Add(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

// Calcu_ServiceDesc is the grpc.ServiceDesc for Calcu service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calcu_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Calcu",
	HandlerType: (*CalcuServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Calcu_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/add.proto",
}
