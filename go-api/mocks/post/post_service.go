package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) Save(ctx context.Context, post *domain.Post) error {

	ret := m.Called(ctx, post)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *MockPostService) FindAll(ctx context.Context) ([]domain.Post, error) {

	ret := m.Called(ctx)

	var r0 []domain.Post
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Post); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}

func (m *MockPostService) FindByID(ctx context.Context, id int32) (domain.Post, error) {
	ret := m.Called(ctx, id)

	var r0 domain.Post
	if rf, ok := ret.Get(0).(func(context.Context, int32) domain.Post); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}
	return r0, r1
}

func (m *MockPostService) Update(ctx context.Context, id int32, post *domain.Post) error {
	ret := m.Called(ctx, id, post)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, *domain.Post) error); ok {
		r0 = rf(ctx, id, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(error)
		}
	}
	return r0
}

func (m *MockPostService) Delete(ctx context.Context, id int32) error {
	ret := m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) error); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(error)
		}
	}
	return r0
}
