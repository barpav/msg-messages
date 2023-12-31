// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/barpav/msg-messages/internal/rest/models"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// CreateNewPersonalMessageV1 provides a mock function with given fields: ctx, sender, data
func (_m *Storage) CreateNewPersonalMessageV1(ctx context.Context, sender string, data *models.NewPersonalMessageV1) (int64, int64, error) {
	ret := _m.Called(ctx, sender, data)

	var r0 int64
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *models.NewPersonalMessageV1) (int64, int64, error)); ok {
		return rf(ctx, sender, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *models.NewPersonalMessageV1) int64); ok {
		r0 = rf(ctx, sender, data)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *models.NewPersonalMessageV1) int64); ok {
		r1 = rf(ctx, sender, data)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, *models.NewPersonalMessageV1) error); ok {
		r2 = rf(ctx, sender, data)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteMessageData provides a mock function with given fields: ctx, id, timestamp
func (_m *Storage) DeleteMessageData(ctx context.Context, id int64, timestamp int64) (int64, error) {
	ret := _m.Called(ctx, id, timestamp)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (int64, error)); ok {
		return rf(ctx, id, timestamp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) int64); ok {
		r0 = rf(ctx, id, timestamp)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, id, timestamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EditMessageText provides a mock function with given fields: ctx, id, timestamp, text
func (_m *Storage) EditMessageText(ctx context.Context, id int64, timestamp int64, text string) (int64, error) {
	ret := _m.Called(ctx, id, timestamp, text)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string) (int64, error)); ok {
		return rf(ctx, id, timestamp, text)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string) int64); ok {
		r0 = rf(ctx, id, timestamp, text)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, string) error); ok {
		r1 = rf(ctx, id, timestamp, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageUpdatesV1 provides a mock function with given fields: ctx, userId, after, limit
func (_m *Storage) MessageUpdatesV1(ctx context.Context, userId string, after int64, limit int) (*models.MessageUpdatesV1, error) {
	ret := _m.Called(ctx, userId, after, limit)

	var r0 *models.MessageUpdatesV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int) (*models.MessageUpdatesV1, error)); ok {
		return rf(ctx, userId, after, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int) *models.MessageUpdatesV1); ok {
		r0 = rf(ctx, userId, after, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MessageUpdatesV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64, int) error); ok {
		r1 = rf(ctx, userId, after, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersonalMessageV1 provides a mock function with given fields: ctx, userId, messageId
func (_m *Storage) PersonalMessageV1(ctx context.Context, userId string, messageId int64) (*models.PersonalMessageV1, error) {
	ret := _m.Called(ctx, userId, messageId)

	var r0 *models.PersonalMessageV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (*models.PersonalMessageV1, error)); ok {
		return rf(ctx, userId, messageId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) *models.PersonalMessageV1); ok {
		r0 = rf(ctx, userId, messageId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PersonalMessageV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, userId, messageId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetMessageReadState provides a mock function with given fields: ctx, id, timestamp, read
func (_m *Storage) SetMessageReadState(ctx context.Context, id int64, timestamp int64, read bool) (int64, error) {
	ret := _m.Called(ctx, id, timestamp, read)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) (int64, error)); ok {
		return rf(ctx, id, timestamp, read)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) int64); ok {
		r0 = rf(ctx, id, timestamp, read)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, bool) error); ok {
		r1 = rf(ctx, id, timestamp, read)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
