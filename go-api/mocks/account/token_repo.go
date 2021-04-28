package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) SetRefreshToken(ctx context.Context, accID string, accessToken string, expiresIn time.Duration) error {
	ret := m.Called(ctx, accID, accessToken, expiresIn)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Duration) error); ok {
		r0 = rf(ctx, accID, accessToken, expiresIn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(error)
		}
	}
	return r0
}

func (m *MockTokenRepo) DeleteRefreshToken(ctx context.Context, accID string, prevAccessToken string) error {
	ret := m.Called(ctx, accID, prevAccessToken)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, accID, prevAccessToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(error)
		}
	}
	return r0
}
