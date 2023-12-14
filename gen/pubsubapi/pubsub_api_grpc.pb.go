// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: pubsubapi/pubsub_api.proto

package pubsubapi

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

// PubSubClient is the client API for PubSub service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PubSubClient interface {
	// Bidirectional streaming RPC to subscribe to a Topic. The subscription is pull-based. A client can request
	// for more events as it consumes events. This enables a client to handle flow control based on the client's processing speed.
	//
	// Typical flow:
	//  1. Client requests for X number of events via FetchRequest.
	//  2. Server receives request and delivers events until X events are delivered to the client via one or more FetchResponse messages.
	//  3. Client consumes the FetchResponse messages as they come.
	//  4. Client issues new FetchRequest for Y more number of events. This request can
	//     come before the server has delivered the earlier requested X number of events
	//     so the client gets a continuous stream of events if any.
	//
	// If a client requests more events before the server finishes the last
	// requested amount, the server appends the new amount to the current amount of
	// events it still needs to fetch and deliver.
	//
	// A client can subscribe at any point in the stream by providing a replay option in the first FetchRequest.
	// The replay option is honored for the first FetchRequest received from a client. Any subsequent FetchRequests with a
	// new replay option are ignored. A client needs to call the Subscribe RPC again to restart the subscription
	// at a new point in the stream.
	//
	// The first FetchRequest of the stream identifies the topic to subscribe to.
	// If any subsequent FetchRequest provides topic_name, it must match what
	// was provided in the first FetchRequest; otherwise, the RPC returns an error
	// with INVALID_ARGUMENT status.
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (PubSub_SubscribeClient, error)
	// Get the event schema for a topic based on a schema ID.
	GetSchema(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*SchemaInfo, error)
	// Get the topic Information related to the specified topic.
	GetTopic(ctx context.Context, in *TopicRequest, opts ...grpc.CallOption) (*TopicInfo, error)
	// Send a publish request to synchronously publish events to a topic.
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	// Bidirectional Streaming RPC to publish events to the event bus.
	// PublishRequest contains the batch of events to publish.
	//
	// The first PublishRequest of the stream identifies the topic to publish on.
	// If any subsequent PublishRequest provides topic_name, it must match what
	// was provided in the first PublishRequest; otherwise, the RPC returns an error
	// with INVALID_ARGUMENT status.
	//
	// The server returns a PublishResponse for each PublishRequest when publish is
	// complete for the batch. A client does not have to wait for a PublishResponse
	// before sending a new PublishRequest, i.e. multiple publish batches can be queued
	// up, which allows for higher publish rate as a client can asynchronously
	// publish more events while publishes are still in flight on the server side.
	//
	// PublishResponse holds a PublishResult for each event published that indicates success
	// or failure of the publish. A client can then retry the publish as needed before sending
	// more PublishRequests for new events to publish.
	//
	// A client must send a valid publish request with one or more events every 70 seconds to hold on to the stream.
	// Otherwise, the server closes the stream and notifies the client. Once the client is notified of the stream closure,
	// it must make a new PublishStream call to resume publishing.
	PublishStream(ctx context.Context, opts ...grpc.CallOption) (PubSub_PublishStreamClient, error)
}

type pubSubClient struct {
	cc grpc.ClientConnInterface
}

func NewPubSubClient(cc grpc.ClientConnInterface) PubSubClient {
	return &pubSubClient{cc}
}

func (c *pubSubClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (PubSub_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &PubSub_ServiceDesc.Streams[0], "/eventbus.v1.PubSub/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &pubSubSubscribeClient{stream}
	return x, nil
}

type PubSub_SubscribeClient interface {
	Send(*FetchRequest) error
	Recv() (*FetchResponse, error)
	grpc.ClientStream
}

type pubSubSubscribeClient struct {
	grpc.ClientStream
}

func (x *pubSubSubscribeClient) Send(m *FetchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pubSubSubscribeClient) Recv() (*FetchResponse, error) {
	m := new(FetchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pubSubClient) GetSchema(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*SchemaInfo, error) {
	out := new(SchemaInfo)
	err := c.cc.Invoke(ctx, "/eventbus.v1.PubSub/GetSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubClient) GetTopic(ctx context.Context, in *TopicRequest, opts ...grpc.CallOption) (*TopicInfo, error) {
	out := new(TopicInfo)
	err := c.cc.Invoke(ctx, "/eventbus.v1.PubSub/GetTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/eventbus.v1.PubSub/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubSubClient) PublishStream(ctx context.Context, opts ...grpc.CallOption) (PubSub_PublishStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PubSub_ServiceDesc.Streams[1], "/eventbus.v1.PubSub/PublishStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &pubSubPublishStreamClient{stream}
	return x, nil
}

type PubSub_PublishStreamClient interface {
	Send(*PublishRequest) error
	Recv() (*PublishResponse, error)
	grpc.ClientStream
}

type pubSubPublishStreamClient struct {
	grpc.ClientStream
}

func (x *pubSubPublishStreamClient) Send(m *PublishRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pubSubPublishStreamClient) Recv() (*PublishResponse, error) {
	m := new(PublishResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PubSubServer is the server API for PubSub service.
// All implementations must embed UnimplementedPubSubServer
// for forward compatibility
type PubSubServer interface {
	// Bidirectional streaming RPC to subscribe to a Topic. The subscription is pull-based. A client can request
	// for more events as it consumes events. This enables a client to handle flow control based on the client's processing speed.
	//
	// Typical flow:
	//  1. Client requests for X number of events via FetchRequest.
	//  2. Server receives request and delivers events until X events are delivered to the client via one or more FetchResponse messages.
	//  3. Client consumes the FetchResponse messages as they come.
	//  4. Client issues new FetchRequest for Y more number of events. This request can
	//     come before the server has delivered the earlier requested X number of events
	//     so the client gets a continuous stream of events if any.
	//
	// If a client requests more events before the server finishes the last
	// requested amount, the server appends the new amount to the current amount of
	// events it still needs to fetch and deliver.
	//
	// A client can subscribe at any point in the stream by providing a replay option in the first FetchRequest.
	// The replay option is honored for the first FetchRequest received from a client. Any subsequent FetchRequests with a
	// new replay option are ignored. A client needs to call the Subscribe RPC again to restart the subscription
	// at a new point in the stream.
	//
	// The first FetchRequest of the stream identifies the topic to subscribe to.
	// If any subsequent FetchRequest provides topic_name, it must match what
	// was provided in the first FetchRequest; otherwise, the RPC returns an error
	// with INVALID_ARGUMENT status.
	Subscribe(PubSub_SubscribeServer) error
	// Get the event schema for a topic based on a schema ID.
	GetSchema(context.Context, *SchemaRequest) (*SchemaInfo, error)
	// Get the topic Information related to the specified topic.
	GetTopic(context.Context, *TopicRequest) (*TopicInfo, error)
	// Send a publish request to synchronously publish events to a topic.
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	// Bidirectional Streaming RPC to publish events to the event bus.
	// PublishRequest contains the batch of events to publish.
	//
	// The first PublishRequest of the stream identifies the topic to publish on.
	// If any subsequent PublishRequest provides topic_name, it must match what
	// was provided in the first PublishRequest; otherwise, the RPC returns an error
	// with INVALID_ARGUMENT status.
	//
	// The server returns a PublishResponse for each PublishRequest when publish is
	// complete for the batch. A client does not have to wait for a PublishResponse
	// before sending a new PublishRequest, i.e. multiple publish batches can be queued
	// up, which allows for higher publish rate as a client can asynchronously
	// publish more events while publishes are still in flight on the server side.
	//
	// PublishResponse holds a PublishResult for each event published that indicates success
	// or failure of the publish. A client can then retry the publish as needed before sending
	// more PublishRequests for new events to publish.
	//
	// A client must send a valid publish request with one or more events every 70 seconds to hold on to the stream.
	// Otherwise, the server closes the stream and notifies the client. Once the client is notified of the stream closure,
	// it must make a new PublishStream call to resume publishing.
	PublishStream(PubSub_PublishStreamServer) error
	mustEmbedUnimplementedPubSubServer()
}

// UnimplementedPubSubServer must be embedded to have forward compatible implementations.
type UnimplementedPubSubServer struct {
}

func (UnimplementedPubSubServer) Subscribe(PubSub_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPubSubServer) GetSchema(context.Context, *SchemaRequest) (*SchemaInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchema not implemented")
}
func (UnimplementedPubSubServer) GetTopic(context.Context, *TopicRequest) (*TopicInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopic not implemented")
}
func (UnimplementedPubSubServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedPubSubServer) PublishStream(PubSub_PublishStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PublishStream not implemented")
}
func (UnimplementedPubSubServer) mustEmbedUnimplementedPubSubServer() {}

// UnsafePubSubServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PubSubServer will
// result in compilation errors.
type UnsafePubSubServer interface {
	mustEmbedUnimplementedPubSubServer()
}

func RegisterPubSubServer(s grpc.ServiceRegistrar, srv PubSubServer) {
	s.RegisterService(&PubSub_ServiceDesc, srv)
}

func _PubSub_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PubSubServer).Subscribe(&pubSubSubscribeServer{stream})
}

type PubSub_SubscribeServer interface {
	Send(*FetchResponse) error
	Recv() (*FetchRequest, error)
	grpc.ServerStream
}

type pubSubSubscribeServer struct {
	grpc.ServerStream
}

func (x *pubSubSubscribeServer) Send(m *FetchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pubSubSubscribeServer) Recv() (*FetchRequest, error) {
	m := new(FetchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PubSub_GetSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServer).GetSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventbus.v1.PubSub/GetSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServer).GetSchema(ctx, req.(*SchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSub_GetTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServer).GetTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventbus.v1.PubSub/GetTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServer).GetTopic(ctx, req.(*TopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSub_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubSubServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventbus.v1.PubSub/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubSubServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PubSub_PublishStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PubSubServer).PublishStream(&pubSubPublishStreamServer{stream})
}

type PubSub_PublishStreamServer interface {
	Send(*PublishResponse) error
	Recv() (*PublishRequest, error)
	grpc.ServerStream
}

type pubSubPublishStreamServer struct {
	grpc.ServerStream
}

func (x *pubSubPublishStreamServer) Send(m *PublishResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pubSubPublishStreamServer) Recv() (*PublishRequest, error) {
	m := new(PublishRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PubSub_ServiceDesc is the grpc.ServiceDesc for PubSub service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PubSub_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "eventbus.v1.PubSub",
	HandlerType: (*PubSubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSchema",
			Handler:    _PubSub_GetSchema_Handler,
		},
		{
			MethodName: "GetTopic",
			Handler:    _PubSub_GetTopic_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _PubSub_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _PubSub_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "PublishStream",
			Handler:       _PubSub_PublishStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pubsubapi/pubsub_api.proto",
}
