package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	mocks "github.com/whuangz/go-example/go-api/mocks/account"
	"github.com/whuangz/go-example/go-api/service"
)

func TestFindMe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockAccResp := &domain.Account{
			UID:   uid,
			Email: "whuangz@gmail.com",
			Name:  "William",
		}

		mockAccountRepository := new(mocks.MockAccountRepo)
		us := service.NewAccountService(mockAccountRepository)

		mockAccountRepository.On("FindByID", mock.Anything, uid).Return(mockAccResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockAccResp)
		mockAccountRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockAccountRepository := new(mocks.MockAccountRepo)
		us := service.NewAccountService(mockAccountRepository)

		mockAccountRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockAccountRepository.AssertExpectations(t)
	})
}
