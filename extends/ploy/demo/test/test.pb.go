// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

/*
Package test is a generated protocol buffer package.

It is generated from these files:
	test.proto

It has these top-level messages:
	TestRequest
	TestResponse
*/
package test

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import (
	ketty "github.com/yyzybb537/ketty"
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

type TestRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *TestRequest) Reset()                    { *m = TestRequest{} }
func (m *TestRequest) String() string            { return proto.CompactTextString(m) }
func (*TestRequest) ProtoMessage()               {}
func (*TestRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TestRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type TestResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *TestResponse) Reset()                    { *m = TestResponse{} }
func (m *TestResponse) String() string            { return proto.CompactTextString(m) }
func (*TestResponse) ProtoMessage()               {}
func (*TestResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TestResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*TestRequest)(nil), "test.TestRequest")
	proto.RegisterType((*TestResponse)(nil), "test.TestResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Test service

type TestClient interface {
	Ping(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error)
}

type testClient struct {
	cc *grpc.ClientConn
}

func NewTestClient(cc *grpc.ClientConn) TestClient {
	return &testClient{cc}
}

func (c *testClient) Ping(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error) {
	out := new(TestResponse)
	err := grpc.Invoke(ctx, "/test.Test/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestServer interface {
	Ping(context.Context, *TestRequest) (*TestResponse, error)
}

func RegisterTestServer(s *grpc.Server, srv TestServer) {
	s.RegisterService(&_Test_serviceDesc, srv)
}

func _Test_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.Test/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).Ping(ctx, req.(*TestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Test_serviceDesc = grpc.ServiceDesc{
	ServiceName: "test.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Test_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test.proto",
}

// Reference imports to suppress errors if they are not otherwise used.
var _ ketty.Dummy

// This is a compile-time assertion to ensure that this generated file
// is compatible with the ketty package it is being compiled against.

type TestHandleT struct {
	desc *grpc.ServiceDesc
}

func (h *TestHandleT) Implement() interface{} {
	return h.desc
}

func (h *TestHandleT) ServiceName() string {
	return h.desc.ServiceName
}

var TestHandle = &TestHandleT{desc: &_Test_serviceDesc}

type KettyTestClient struct {
	client ketty.Client
}

func NewKettyTestClient(client ketty.Client) *KettyTestClient {
	return &KettyTestClient{client}
}

func (this *KettyTestClient) Ping(ctx context.Context, in *TestRequest) (*TestResponse, error) {
	out := new(TestResponse)
	err := this.client.Invoke(ctx, TestHandle, "Ping", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func init() { proto.RegisterFile("test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x14, 0xb9, 0xb8, 0x43, 0x52,
	0x8b, 0x4b, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73,
	0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x25, 0x0d, 0x2e, 0x1e, 0x88, 0x92,
	0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x09, 0x2e, 0xf6, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0x74,
	0x98, 0x32, 0x18, 0xd7, 0xc8, 0x9c, 0x8b, 0x05, 0xa4, 0x52, 0x48, 0x9f, 0x8b, 0x25, 0x20, 0x33,
	0x2f, 0x5d, 0x48, 0x50, 0x0f, 0x6c, 0x1f, 0x92, 0x05, 0x52, 0x42, 0xc8, 0x42, 0x10, 0x03, 0x95,
	0x18, 0x92, 0xd8, 0xc0, 0x4e, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xdf, 0xc3, 0x18, 0x23,
	0xa0, 0x00, 0x00, 0x00,
}
