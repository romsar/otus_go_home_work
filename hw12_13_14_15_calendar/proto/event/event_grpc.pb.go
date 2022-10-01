// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: event/event.proto

package event

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

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	CreateEventV1(ctx context.Context, in *CreateEventRequestV1, opts ...grpc.CallOption) (*EventReplyV1, error)
	UpdateEventV1(ctx context.Context, in *UpdateEventRequestV1, opts ...grpc.CallOption) (*EventReplyV1, error)
	DeleteEventV1(ctx context.Context, in *DeleteEventRequestV1, opts ...grpc.CallOption) (*DeleteEventReplyV1, error)
	GetEventsForDayV1(ctx context.Context, in *GetEventsForDayRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error)
	GetEventsForWeekV1(ctx context.Context, in *GetEventsForWeekRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error)
	GetEventsForMonthV1(ctx context.Context, in *GetEventsForMonthRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEventV1(ctx context.Context, in *CreateEventRequestV1, opts ...grpc.CallOption) (*EventReplyV1, error) {
	out := new(EventReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/CreateEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEventV1(ctx context.Context, in *UpdateEventRequestV1, opts ...grpc.CallOption) (*EventReplyV1, error) {
	out := new(EventReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/UpdateEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEventV1(ctx context.Context, in *DeleteEventRequestV1, opts ...grpc.CallOption) (*DeleteEventReplyV1, error) {
	out := new(DeleteEventReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/DeleteEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEventsForDayV1(ctx context.Context, in *GetEventsForDayRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error) {
	out := new(EventsReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/GetEventsForDayV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEventsForWeekV1(ctx context.Context, in *GetEventsForWeekRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error) {
	out := new(EventsReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/GetEventsForWeekV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEventsForMonthV1(ctx context.Context, in *GetEventsForMonthRequestV1, opts ...grpc.CallOption) (*EventsReplyV1, error) {
	out := new(EventsReplyV1)
	err := c.cc.Invoke(ctx, "/event.EventService/GetEventsForMonthV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	CreateEventV1(context.Context, *CreateEventRequestV1) (*EventReplyV1, error)
	UpdateEventV1(context.Context, *UpdateEventRequestV1) (*EventReplyV1, error)
	DeleteEventV1(context.Context, *DeleteEventRequestV1) (*DeleteEventReplyV1, error)
	GetEventsForDayV1(context.Context, *GetEventsForDayRequestV1) (*EventsReplyV1, error)
	GetEventsForWeekV1(context.Context, *GetEventsForWeekRequestV1) (*EventsReplyV1, error)
	GetEventsForMonthV1(context.Context, *GetEventsForMonthRequestV1) (*EventsReplyV1, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) CreateEventV1(context.Context, *CreateEventRequestV1) (*EventReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEventV1 not implemented")
}
func (UnimplementedEventServiceServer) UpdateEventV1(context.Context, *UpdateEventRequestV1) (*EventReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEventV1 not implemented")
}
func (UnimplementedEventServiceServer) DeleteEventV1(context.Context, *DeleteEventRequestV1) (*DeleteEventReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEventV1 not implemented")
}
func (UnimplementedEventServiceServer) GetEventsForDayV1(context.Context, *GetEventsForDayRequestV1) (*EventsReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForDayV1 not implemented")
}
func (UnimplementedEventServiceServer) GetEventsForWeekV1(context.Context, *GetEventsForWeekRequestV1) (*EventsReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForWeekV1 not implemented")
}
func (UnimplementedEventServiceServer) GetEventsForMonthV1(context.Context, *GetEventsForMonthRequestV1) (*EventsReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForMonthV1 not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_CreateEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/CreateEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEventV1(ctx, req.(*CreateEventRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/UpdateEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEventV1(ctx, req.(*UpdateEventRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/DeleteEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEventV1(ctx, req.(*DeleteEventRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEventsForDayV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsForDayRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEventsForDayV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEventsForDayV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEventsForDayV1(ctx, req.(*GetEventsForDayRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEventsForWeekV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsForWeekRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEventsForWeekV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEventsForWeekV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEventsForWeekV1(ctx, req.(*GetEventsForWeekRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEventsForMonthV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsForMonthRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEventsForMonthV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEventsForMonthV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEventsForMonthV1(ctx, req.(*GetEventsForMonthRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEventV1",
			Handler:    _EventService_CreateEventV1_Handler,
		},
		{
			MethodName: "UpdateEventV1",
			Handler:    _EventService_UpdateEventV1_Handler,
		},
		{
			MethodName: "DeleteEventV1",
			Handler:    _EventService_DeleteEventV1_Handler,
		},
		{
			MethodName: "GetEventsForDayV1",
			Handler:    _EventService_GetEventsForDayV1_Handler,
		},
		{
			MethodName: "GetEventsForWeekV1",
			Handler:    _EventService_GetEventsForWeekV1_Handler,
		},
		{
			MethodName: "GetEventsForMonthV1",
			Handler:    _EventService_GetEventsForMonthV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event/event.proto",
}