// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package message

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

// MessageClient is the client API for Message service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageClient interface {
	Push(ctx context.Context, opts ...grpc.CallOption) (Message_PushClient, error)
}

type messageClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageClient(cc grpc.ClientConnInterface) MessageClient {
	return &messageClient{cc}
}

func (c *messageClient) Push(ctx context.Context, opts ...grpc.CallOption) (Message_PushClient, error) {
	stream, err := c.cc.NewStream(ctx, &Message_ServiceDesc.Streams[0], "/message.message/Push", opts...)
	if err != nil {
		return nil, err
	}
	x := &messagePushClient{stream}
	return x, nil
}

type Message_PushClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type messagePushClient struct {
	grpc.ClientStream
}

func (x *messagePushClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messagePushClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServer is the server API for Message service.
// All implementations must embed UnimplementedMessageServer
// for forward compatibility
type MessageServer interface {
	Push(Message_PushServer) error
	mustEmbedUnimplementedMessageServer()
}

// UnimplementedMessageServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServer struct {
}

func (UnimplementedMessageServer) Push(Message_PushServer) error {
	return status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedMessageServer) mustEmbedUnimplementedMessageServer() {}

// UnsafeMessageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServer will
// result in compilation errors.
type UnsafeMessageServer interface {
	mustEmbedUnimplementedMessageServer()
}

func RegisterMessageServer(s grpc.ServiceRegistrar, srv MessageServer) {
	s.RegisterService(&Message_ServiceDesc, srv)
}

func _Message_Push_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageServer).Push(&messagePushServer{stream})
}

type Message_PushServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type messagePushServer struct {
	grpc.ServerStream
}

func (x *messagePushServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messagePushServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Message_ServiceDesc is the grpc.ServiceDesc for Message service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Message_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.message",
	HandlerType: (*MessageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Push",
			Handler:       _Message_Push_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "message/message.proto",
}
