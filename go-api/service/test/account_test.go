package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/helpers/crypto"
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
		a, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, a, mockAccResp)
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

func TestSignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockAcc := &domain.Account{
			Email:    "whuangz@gmail.com",
			Password: "admin123",
		}

		mockAccountRepository := new(mocks.MockAccountRepo)
		us := service.NewAccountService(mockAccountRepository)

		mockAccountRepository.On("Create", mock.Anything, mockAcc).
			Run(func(args mock.Arguments) {
				accArg := args.Get(1).(*domain.Account)
				accArg.UID = uid
			}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockAcc)

		assert.NoError(t, err)
		assert.Equal(t, uid, mockAcc.UID)
		mockAccountRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockAcc := &domain.Account{
			Email:    "whuangz@gmail.com",
			Password: "admin123",
		}

		mockAccountRepository := new(mocks.MockAccountRepo)
		us := service.NewAccountService(mockAccountRepository)

		mockErr := domain.NewConflict("email", mockAcc.Email)

		mockAccountRepository.On("Create", mock.Anything, mockAcc).Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockAcc)

		assert.EqualError(t, err, mockErr.Error())

		mockAccountRepository.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	email := "whuangz@gmail.com"
	validPW := "admin123"
	hashedValidPW, _ := crypto.HashPassword(validPW)
	invalidPW := "howdyhodufus!"

	mockAccRepo := new(mocks.MockAccountRepo)
	us := service.NewAccountService(mockAccRepo)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &domain.Account{
			Email:    email,
			Password: validPW,
		}

		mockUserResp := &domain.Account{
			UID:      uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockAccRepo.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.NoError(t, err)
		mockAccRepo.AssertCalled(t, "FindByEmail", mockArgs...)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &domain.Account{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &domain.Account{
			UID:      uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockAccRepo.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, "Invalid email and password combination")
		mockAccRepo.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}
