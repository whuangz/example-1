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
