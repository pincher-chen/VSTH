// Code generated by protoc-gen-go.
// source: addsvc.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	addsvc.proto

It has these top-level messages:
	SumRequest
	SumReply
	ConcatRequest
	ConcatReply
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "github.com/golang/net/context"
	grpc "github.com/grpc/grpc-go"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The sum request contains two parameters.
type SumRequest struct {
	A int64 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
	B int64 `protobuf:"varint,2,opt,name=b" json:"b,omitempty"`
}

func (m *SumRequest) Reset()                    { *m = SumRequest{} }
func (m *SumRequest) String() string            { return proto.CompactTextString(m) }
func (*SumRequest) ProtoMessage()               {}
func (*SumRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The sum response contains the result of the calculation.
type SumReply struct {
	V   int64  `protobuf:"varint,1,opt,name=v" json:"v,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *SumReply) Reset()                    { *m = SumReply{} }
func (m *SumReply) String() string            { return proto.CompactTextString(m) }
func (*SumReply) ProtoMessage()               {}
func (*SumReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// The Concat request contains two parameters.
type ConcatRequest struct {
	A string `protobuf:"bytes,1,opt,name=a" json:"a,omitempty"`
	B string `protobuf:"bytes,2,opt,name=b" json:"b,omitempty"`
}

func (m *ConcatRequest) Reset()                    { *m = ConcatRequest{} }
func (m *ConcatRequest) String() string            { return proto.CompactTextString(m) }
func (*ConcatRequest) ProtoMessage()               {}
func (*ConcatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// The Concat response contains the result of the concatenation.
type ConcatReply struct {
	V   string `protobuf:"bytes,1,opt,name=v" json:"v,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *ConcatReply) Reset()                    { *m = ConcatReply{} }
func (m *ConcatReply) String() string            { return proto.CompactTextString(m) }
func (*ConcatReply) ProtoMessage()               {}
func (*ConcatReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*SumRequest)(nil), "pb.SumRequest")
	proto.RegisterType((*SumReply)(nil), "pb.SumReply")
	proto.RegisterType((*ConcatRequest)(nil), "pb.ConcatRequest")
	proto.RegisterType((*ConcatReply)(nil), "pb.ConcatReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Add service

type AddClient interface {
	// Sums two integers.
	Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumReply, error)
	// Concatenates two strings
	Concat(ctx context.Context, in *ConcatRequest, opts ...grpc.CallOption) (*ConcatReply, error)
}

type addClient struct {
	cc *grpc.ClientConn
}

func NewAddClient(cc *grpc.ClientConn) AddClient {
	return &addClient{cc}
}

func (c *addClient) Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumReply, error) {
	out := new(SumReply)
	err := grpc.Invoke(ctx, "/pb.Add/Sum", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addClient) Concat(ctx context.Context, in *ConcatRequest, opts ...grpc.CallOption) (*ConcatReply, error) {
	out := new(ConcatReply)
	err := grpc.Invoke(ctx, "/pb.Add/Concat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Add service

type AddServer interface {
	// Sums two integers.
	Sum(context.Context, *SumRequest) (*SumReply, error)
	// Concatenates two strings
	Concat(context.Context, *ConcatRequest) (*ConcatReply, error)
}

func RegisterAddServer(s *grpc.Server, srv AddServer) {
	s.RegisterService(&_Add_serviceDesc, srv)
}

func _Add_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Add/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServer).Sum(ctx, req.(*SumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Add_Concat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConcatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServer).Concat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Add/Concat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServer).Concat(ctx, req.(*ConcatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Add_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Add",
	HandlerType: (*AddServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _Add_Sum_Handler,
		},
		{
			MethodName: "Concat",
			Handler:    _Add_Concat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("addsvc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x4c, 0x49, 0x29,
	0x2e, 0x4b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xd2, 0xe0, 0xe2,
	0x0a, 0x2e, 0xcd, 0x0d, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0xe2, 0xe1, 0x62, 0x4c, 0x94,
	0x60, 0x54, 0x60, 0xd4, 0x60, 0x0e, 0x62, 0x4c, 0x04, 0xf1, 0x92, 0x24, 0x98, 0x20, 0xbc, 0x24,
	0x25, 0x2d, 0x2e, 0x0e, 0xb0, 0xca, 0x82, 0x9c, 0x4a, 0x90, 0x4c, 0x19, 0x4c, 0x5d, 0x99, 0x90,
	0x00, 0x17, 0x73, 0x6a, 0x51, 0x11, 0x58, 0x25, 0x67, 0x10, 0x88, 0xa9, 0xa4, 0xcd, 0xc5, 0xeb,
	0x9c, 0x9f, 0x97, 0x9c, 0x58, 0x82, 0x61, 0x30, 0x27, 0x8a, 0xc1, 0x9c, 0x20, 0x83, 0x75, 0xb9,
	0xb8, 0x61, 0x8a, 0x51, 0xcc, 0xe6, 0xc4, 0x6a, 0xb6, 0x51, 0x0c, 0x17, 0xb3, 0x63, 0x4a, 0x8a,
	0x90, 0x2a, 0x17, 0x33, 0xd0, 0x39, 0x42, 0x7c, 0x7a, 0x05, 0x49, 0x7a, 0x08, 0x1f, 0x48, 0xf1,
	0xc0, 0xf9, 0x40, 0xb3, 0x94, 0x18, 0x84, 0xf4, 0xb8, 0xd8, 0x20, 0x86, 0x0b, 0x09, 0x82, 0x64,
	0x50, 0x5c, 0x25, 0xc5, 0x8f, 0x2c, 0x04, 0x56, 0x9f, 0xc4, 0x06, 0x0e, 0x1a, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xdc, 0x37, 0x81, 0x99, 0x2a, 0x01, 0x00, 0x00,
}
