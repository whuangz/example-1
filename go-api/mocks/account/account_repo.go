package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
)

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	ret := m.Called(ctx, email)
	var r0 *domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Account); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}

func (m *MockAccountRepo) FindByID(ctx context.Context, uid uuid.UUID) (*domain.Account, error) {

	ret := m.Called(ctx, uid)

	var r0 *domain.Account
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.Account); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			// we can just return this if we know we won't be passing function to "Return"
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, uid)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}

func (m *MockAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	ret := m.Called(ctx, a)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Account) error); ok {
		r0 = rf(ctx, a)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(error)
		}
	}
	return r0
}
