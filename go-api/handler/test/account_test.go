package handle_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/handler"
	mocks "github.com/whuangz/go-example/go-api/mocks/account"
)

func TestGetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)

		uid, _ := uuid.NewRandom()

		mockAccResp := &domain.Account{
			UID:   uid,
			Email: "Whuangz@gmail.com",
			Name:  "William",
		}

		mockAccService.On("Get", mock.AnythingOfType("*context.emptyCtx"), uid).Return(mockAccResp, nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("account", &domain.Account{
				UID: uid,
			},
			)
		})

		handler.NewAccountHandler(router, mockAccService, nil)
		request, err := http.NewRequest(http.MethodGet, "/api/account/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, request)

		respBody, err := json.Marshal(gin.H{
			"data": mockAccResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, respBody, rec.Body.Bytes())
		mockAccService.AssertExpectations(t)
	})

	t.Run("No context found", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)

		mockAccService.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, nil)
		request, err := http.NewRequest(http.MethodGet, "/api/account/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, request)

		assert.Equal(t, 401, rec.Code)

		mockAccService.AssertNotCalled(t, "Get", mock.Anything)

	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mockAccService := new(mocks.MockAccountService)
		mockAccService.On("Get", mock.AnythingOfType("*context.emptyCtx"), uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		// a response recorder for getting written http response
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("account", &domain.Account{
				UID: uid,
			},
			)
		})

		handler.NewAccountHandler(router, mockAccService, nil)

		request, err := http.NewRequest(http.MethodGet, "/api/account/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, request)

		respErr := domain.NewNotFound("account", uid.String())

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})
		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), rec.Code)
		assert.Equal(t, respBody, rec.Body.Bytes())
		mockAccService.AssertExpectations(t)
	})
}

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Email and Password Required", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)
		mockAccService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*domain.Account")).Return(nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, nil)

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, request)

		assert.Equal(t, 400, rec.Code)
		mockAccService.AssertNotCalled(t, "Signup")
	})

	t.Run("Invalid Email", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)
		mockAccService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*domain.Account")).Return(nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, nil)

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "&*&@gmail",
			"password": "avalidpassword123",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, request)

		assert.Equal(t, 400, rec.Code)
		mockAccService.AssertNotCalled(t, "Signup")
	})

	t.Run("Error calling UserService", func(t *testing.T) {
		a := &domain.Account{
			Email:    "whuangz@gmail.com",
			Password: "avalidpassword",
		}

		mockAccService := new(mocks.MockAccountService)
		mockAccService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), a).Return(domain.NewConflict("User Already Exists", a.Email))

		// a response recorder for getting written http response
		rec := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, nil)

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    a.Email,
			"password": a.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, request)

		assert.Equal(t, 409, rec.Code)
		mockAccService.AssertExpectations(t)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)
		mockTokenService := new(mocks.MockTokenService)

		a := &domain.Account{
			Email:    "Whuangz@gmail.com",
			Password: "admin123",
		}

		mockTokenResp := &domain.TokenPair{
			AccessToken:  "idToken",
			RefreshToken: "refreshToken",
		}

		mockAccService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), a).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), a, "").Return(mockTokenResp, nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, mockTokenService)

		reqBody, err := json.Marshal(gin.H{
			"email":    a.Email,
			"password": a.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, request)

		respBody, err := json.Marshal(gin.H{
			"data": mockTokenResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, respBody, rec.Body.Bytes())
		mockAccService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)

	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		mockAccService := new(mocks.MockAccountService)
		mockTokenService := new(mocks.MockTokenService)

		a := &domain.Account{
			Email:    "Whuangz@gmail.com",
			Password: "admin123",
		}

		mockErrorResponse := domain.NewInternal()

		mockAccService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), a).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), a, "").Return(nil, mockErrorResponse)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService, mockTokenService)

		reqBody, err := json.Marshal(gin.H{
			"email":    a.Email,
			"password": a.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rec, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorResponse.Status(), rec.Code)
		assert.Equal(t, respBody, rec.Body.Bytes())
		mockAccService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockAccService := new(mocks.MockAccountService)
	mockTokenService := new(mocks.MockTokenService)

	router := gin.Default()
	handler.NewAccountHandler(router, mockAccService, mockTokenService)

	t.Run("Bad request data", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with invalid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    "notanemail",
			"password": "short",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockAccService.AssertNotCalled(t, "Signin")
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
	})

	t.Run("Error Returned from UserService.Signin", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwdoesnotmatch123"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&domain.Account{Email: email, Password: password},
		}

		// so we can check for a known status code
		mockError := domain.NewAuthorization("invalid email/password combo")

		mockAccService.On("Signin", mockUSArgs...).Return(mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockAccService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		email := "whuangz@gmail.com"
		password := "admin123"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&domain.Account{Email: email, Password: password},
		}

		mockAccService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&domain.Account{Email: email, Password: password},
			"",
		}

		mockTokenPair := &domain.TokenPair{
			AccessToken:  "acctoken",
			RefreshToken: "refreshToken",
		}

		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(mockTokenPair, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"data": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockAccService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		email := "cannotproducetoken@bob.com"
		password := "cannotproducetoken"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&domain.Account{Email: email, Password: password},
		}

		mockAccService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&domain.Account{Email: email, Password: password},
			"",
		}

		mockError := domain.NewInternal()
		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(nil, mockError)
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockAccService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})
}
