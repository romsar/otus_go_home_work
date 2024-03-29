// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	event "github.com/RomanSarvarov/otus_go_home_work/calendar/proto/event"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// EventServiceClient is an autogenerated mock type for the EventServiceClient type
type EventServiceClient struct {
	mock.Mock
}

// CreateEventV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) CreateEventV1(ctx context.Context, in *event.CreateEventRequestV1, opts ...grpc.CallOption) (*event.EventResponseV1, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *event.EventResponseV1
	if rf, ok := ret.Get(0).(func(context.Context, *event.CreateEventRequestV1, ...grpc.CallOption) *event.EventResponseV1); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.EventResponseV1)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.CreateEventRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEventV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) DeleteEventV1(ctx context.Context, in *event.DeleteEventRequestV1, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *event.DeleteEventRequestV1, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.DeleteEventRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForDayV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) GetEventsForDayV1(ctx context.Context, in *event.GetEventsForDayRequestV1, opts ...grpc.CallOption) (*event.EventsResponseV1, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *event.EventsResponseV1
	if rf, ok := ret.Get(0).(func(context.Context, *event.GetEventsForDayRequestV1, ...grpc.CallOption) *event.EventsResponseV1); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.EventsResponseV1)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.GetEventsForDayRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForMonthV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) GetEventsForMonthV1(ctx context.Context, in *event.GetEventsForMonthRequestV1, opts ...grpc.CallOption) (*event.EventsResponseV1, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *event.EventsResponseV1
	if rf, ok := ret.Get(0).(func(context.Context, *event.GetEventsForMonthRequestV1, ...grpc.CallOption) *event.EventsResponseV1); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.EventsResponseV1)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.GetEventsForMonthRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForWeekV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) GetEventsForWeekV1(ctx context.Context, in *event.GetEventsForWeekRequestV1, opts ...grpc.CallOption) (*event.EventsResponseV1, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *event.EventsResponseV1
	if rf, ok := ret.Get(0).(func(context.Context, *event.GetEventsForWeekRequestV1, ...grpc.CallOption) *event.EventsResponseV1); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.EventsResponseV1)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.GetEventsForWeekRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEventV1 provides a mock function with given fields: ctx, in, opts
func (_m *EventServiceClient) UpdateEventV1(ctx context.Context, in *event.UpdateEventRequestV1, opts ...grpc.CallOption) (*event.EventResponseV1, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *event.EventResponseV1
	if rf, ok := ret.Get(0).(func(context.Context, *event.UpdateEventRequestV1, ...grpc.CallOption) *event.EventResponseV1); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*event.EventResponseV1)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *event.UpdateEventRequestV1, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEventServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventServiceClient creates a new instance of EventServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventServiceClient(t mockConstructorTestingTNewEventServiceClient) *EventServiceClient {
	mock := &EventServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
