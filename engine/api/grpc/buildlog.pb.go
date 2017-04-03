// Code generated by protoc-gen-go.
// source: buildlog.proto
// DO NOT EDIT!

/*
Package grpc is a generated protocol buffer package.

It is generated from these files:
	buildlog.proto

It has these top-level messages:
*/
package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import github_com_ovh_cds_sdk "github.com/ovh/cds/sdk"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc1 "google.golang.org/grpc"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc1.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc1.SupportPackageIsVersion4

// Client API for BuildLog service

type BuildLogClient interface {
	AddBuildLog(ctx context.Context, opts ...grpc1.CallOption) (BuildLog_AddBuildLogClient, error)
}

type buildLogClient struct {
	cc *grpc1.ClientConn
}

func NewBuildLogClient(cc *grpc1.ClientConn) BuildLogClient {
	return &buildLogClient{cc}
}

func (c *buildLogClient) AddBuildLog(ctx context.Context, opts ...grpc1.CallOption) (BuildLog_AddBuildLogClient, error) {
	stream, err := grpc1.NewClientStream(ctx, &_BuildLog_serviceDesc.Streams[0], c.cc, "/grpc.BuildLog/AddBuildLog", opts...)
	if err != nil {
		return nil, err
	}
	x := &buildLogAddBuildLogClient{stream}
	return x, nil
}

type BuildLog_AddBuildLogClient interface {
	Send(*github_com_ovh_cds_sdk.Log) error
	CloseAndRecv() (*google_protobuf1.Empty, error)
	grpc1.ClientStream
}

type buildLogAddBuildLogClient struct {
	grpc1.ClientStream
}

func (x *buildLogAddBuildLogClient) Send(m *github_com_ovh_cds_sdk.Log) error {
	return x.ClientStream.SendMsg(m)
}

func (x *buildLogAddBuildLogClient) CloseAndRecv() (*google_protobuf1.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(google_protobuf1.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for BuildLog service

type BuildLogServer interface {
	AddBuildLog(BuildLog_AddBuildLogServer) error
}

func RegisterBuildLogServer(s *grpc1.Server, srv BuildLogServer) {
	s.RegisterService(&_BuildLog_serviceDesc, srv)
}

func _BuildLog_AddBuildLog_Handler(srv interface{}, stream grpc1.ServerStream) error {
	return srv.(BuildLogServer).AddBuildLog(&buildLogAddBuildLogServer{stream})
}

type BuildLog_AddBuildLogServer interface {
	SendAndClose(*google_protobuf1.Empty) error
	Recv() (*github_com_ovh_cds_sdk.Log, error)
	grpc1.ServerStream
}

type buildLogAddBuildLogServer struct {
	grpc1.ServerStream
}

func (x *buildLogAddBuildLogServer) SendAndClose(m *google_protobuf1.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *buildLogAddBuildLogServer) Recv() (*github_com_ovh_cds_sdk.Log, error) {
	m := new(github_com_ovh_cds_sdk.Log)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _BuildLog_serviceDesc = grpc1.ServiceDesc{
	ServiceName: "grpc.BuildLog",
	HandlerType: (*BuildLogServer)(nil),
	Methods:     []grpc1.MethodDesc{},
	Streams: []grpc1.StreamDesc{
		{
			StreamName:    "AddBuildLog",
			Handler:       _BuildLog_AddBuildLog_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "buildlog.proto",
}

func init() { proto.RegisterFile("buildlog.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x2a, 0xcd, 0xcc,
	0x49, 0xc9, 0xc9, 0x4f, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a, 0x48,
	0x96, 0x52, 0x48, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x2f, 0xcb,
	0xd0, 0x4f, 0x4e, 0x29, 0xd6, 0x2f, 0x4e, 0xc9, 0xd6, 0x87, 0xab, 0x93, 0x92, 0x4e, 0xcf, 0xcf,
	0x4f, 0xcf, 0x49, 0xd5, 0x07, 0xf3, 0x92, 0x4a, 0xd3, 0xf4, 0x53, 0x73, 0x0b, 0x4a, 0x2a, 0x21,
	0x92, 0x46, 0x41, 0x5c, 0x1c, 0x4e, 0x20, 0x63, 0x7d, 0xf2, 0xd3, 0x85, 0xdc, 0xb8, 0xb8, 0x1d,
	0x53, 0x52, 0xe0, 0x5c, 0x69, 0x3d, 0x84, 0xd1, 0x7a, 0xf9, 0x65, 0x19, 0x7a, 0xc9, 0x29, 0xc5,
	0x7a, 0xc5, 0x29, 0xd9, 0x7a, 0x3e, 0xf9, 0xe9, 0x52, 0x62, 0x7a, 0x10, 0x53, 0xf5, 0x60, 0xa6,
	0xea, 0xb9, 0x82, 0x4c, 0x55, 0x62, 0xd0, 0x60, 0x4c, 0x62, 0x03, 0x8b, 0x19, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xa7, 0xd0, 0xa7, 0x82, 0xb1, 0x00, 0x00, 0x00,
}
