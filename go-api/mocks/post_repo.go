package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
)

type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Save(ctx context.Context, post *domain.Post) error {

	ret := m.Called(ctx, post)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *MockPostRepo) FindAll(ctx context.Context) ([]domain.Post, error) {

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
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *MockPostRepo) FindByID(ctx context.Context, id uint32) (domain.Post, error) {
	return domain.Post{}, nil
}

func (m *MockPostRepo) Update(ctx context.Context, id uint32, post *domain.Post) error {
	return nil
}

func (m *MockPostRepo) Delete(ctx context.Context, id uint32) (int32, error) {
	return 0, nil
}
