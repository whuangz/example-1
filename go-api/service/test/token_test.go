package service_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	jwtHelper "github.com/whuangz/go-example/go-api/helpers/jwt"
	mocks "github.com/whuangz/go-example/go-api/mocks/account"
	"github.com/whuangz/go-example/go-api/service"
)

func TestNewPairFromUser(t *testing.T) {
	var accExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600

	priv, _ := ioutil.ReadFile("../../config/rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../../config/rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	// instantiate a common token service to be used by all tests
	mockTokenRepository := new(mocks.MockTokenRepo)
	tokenService := service.NewTokenService(mockTokenRepository, privKey, pubKey, secret, accExp, refreshExp)

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	a := &domain.Account{
		UID:      uid,
		Email:    "whuangz@gmail.com",
		Password: "admin123",
	}

	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &domain.Account{
		UID:      uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevAccessToken := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		a.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uErrorCase.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		a.UID.String(),
		prevAccessToken,
	}

	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)

	t.Run("Returns a token pair with values", func(t *testing.T) {
		ctx := context.Background()
		tokenPair, err := tokenService.NewPairFromUser(ctx, a, prevAccessToken)
		assert.NoError(t, err)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.AccessToken)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		accTokenClaims := &jwtHelper.AccessTokenTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.AccessToken, accTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			a.UID,
			a.Email,
			a.Name,
			a.ImageUrl,
			a.Website,
		}
		actualIDClaims := []interface{}{
			accTokenClaims.Account.UID,
			accTokenClaims.Account.Email,
			accTokenClaims.Account.Name,
			accTokenClaims.Account.ImageUrl,
			accTokenClaims.Account.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, accTokenClaims.Account.Password) // password should never be encoded to json

		expiresAt := time.Unix(accTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(accExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &jwtHelper.RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, a.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, uErrorCase, "")
		assert.Error(t, err) // should return an error

		// SetRefreshToken should be called with setErrorArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for Prev acc token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, a, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
}

func TestValidateIDToken(t *testing.T) {
	var accExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600

	priv, _ := ioutil.ReadFile("../../config/rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../../config/rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	// instantiate a common token service to be used by all tests
	mockTokenRepository := new(mocks.MockTokenRepo)
	tokenService := service.NewTokenService(mockTokenRepository, privKey, pubKey, secret, accExp, refreshExp)

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	a := &domain.Account{
		UID:      uid,
		Email:    "whuangz@gmail.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := jwtHelper.GenerateAccessToken(a, privKey, accExp)

		aFromToken, err := tokenService.ValidateAccessToken(ss)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{a.Email, a.Name, a.UID, a.Website, a.ImageUrl},
			[]interface{}{aFromToken.Email, aFromToken.Name, aFromToken.UID, aFromToken.Website, aFromToken.ImageUrl},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := jwtHelper.GenerateAccessToken(a, privKey, -1) // expires one second ago

		expectedErr := domain.NewAuthorization("Unable to verify user")

		_, err := tokenService.ValidateAccessToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})
}

func TestValidateRefreshToken(t *testing.T) {
	var accExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600

	priv, _ := ioutil.ReadFile("../../config/rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../../config/rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	// instantiate a common token service to be used by all tests
	mockTokenRepository := new(mocks.MockTokenRepo)
	tokenService := service.NewTokenService(mockTokenRepository, privKey, pubKey, secret, accExp, refreshExp)

	uid, _ := uuid.NewRandom()
	a := &domain.Account{
		UID:      uid,
		Email:    "whuangz@gmail.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {

		testRefreshToken, _ := jwtHelper.GenerateRefreshToken(a.UID, secret, refreshExp)

		validatedRefreshToken, err := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.NoError(t, err)

		assert.Equal(t, a.UID, validatedRefreshToken.UID)
		assert.Equal(t, testRefreshToken.SS, validatedRefreshToken.SS)
		assert.Equal(t, a.UID, validatedRefreshToken.UID)
	})

	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := jwtHelper.GenerateRefreshToken(a.UID, secret, -1)

		expectedErr := domain.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.EqualError(t, err, expectedErr.Message)
	})
}
