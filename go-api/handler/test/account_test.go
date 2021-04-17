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

		mockAccService.On("Get", mock.AnythingOfType("*gin.Context"), uid).Return(mockAccResp, nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("account", &domain.Account{
				UID: uid,
			},
			)
		})

		handler.NewAccountHandler(router, mockAccService)
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
		handler.NewAccountHandler(router, mockAccService)
		request, err := http.NewRequest(http.MethodGet, "/api/account/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, request)

		assert.Equal(t, 401, rec.Code)

		mockAccService.AssertNotCalled(t, "Get", mock.Anything)

	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mockAccService := new(mocks.MockAccountService)
		mockAccService.On("Get", mock.AnythingOfType("*gin.Context"), uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		// a response recorder for getting written http response
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("account", &domain.Account{
				UID: uid,
			},
			)
		})

		handler.NewAccountHandler(router, mockAccService)

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
		mockAccService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*domain.Account")).Return(nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService)

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
		mockAccService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*domain.Account")).Return(nil)

		rec := httptest.NewRecorder()

		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService)

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
		mockAccService.On("Signup", mock.AnythingOfType("*gin.Context"), a).Return(domain.NewConflict("User Already Exists", a.Email))

		// a response recorder for getting written http response
		rec := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()
		handler.NewAccountHandler(router, mockAccService)

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
}
