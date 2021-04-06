package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
)

type MockPostService struct {
	mock.Mock
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
