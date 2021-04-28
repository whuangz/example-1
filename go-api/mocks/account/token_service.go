package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewPairFromUser(ctx context.Context, a *domain.Account, prevAccesstoken string) (*domain.TokenPair, error) {
	ret := m.Called(ctx, a, prevAccesstoken)

	var r0 *domain.TokenPair
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Account, string) *domain.TokenPair); ok {
		r0 = rf(ctx, a, prevAccesstoken)
	} else {
		if ret.Get(0) != nil {
			// we can just return this if we know we won't be passing function to "Return"
			r0 = ret.Get(0).(*domain.TokenPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.Account, string) error); ok {
		r1 = rf(ctx, a, prevAccesstoken)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}

func (m *MockTokenService) ValidateAccessToken(token string) (*domain.Account, error) {
	ret := m.Called(token)

	var r0 *domain.Account
	if rf, ok := ret.Get(0).(func(string) *domain.Account); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			// we can just return this if we know we won't be passing function to "Return"
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}
func (m *MockTokenService) ValidateRefreshToken(tokenString string) (*domain.RefreshToken, error) {
	ret := m.Called(tokenString)

	var r0 *domain.RefreshToken
	if rf, ok := ret.Get(0).(func(string) *domain.RefreshToken); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			// we can just return this if we know we won't be passing function to "Return"
			r0 = ret.Get(0).(*domain.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(error)
		}
	}

	return r0, r1
}
