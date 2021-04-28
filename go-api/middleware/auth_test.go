package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/whuangz/go-example/go-api/domain"
	mocks "github.com/whuangz/go-example/go-api/mocks/account"
)

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)

	uid, _ := uuid.NewRandom()
	u := &domain.Account{
		UID:   uid,
		Email: "whuangz@gmai.com",
	}

	// Since we mock tokenService, we need not
	// create actual JWTs
	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenErr := domain.NewAuthorization("Unable to verify user from Acctoken")

	t.Run("Adds a user to context", func(t *testing.T) {
		mockTokenService.On("ValidateAccessToken", validTokenHeader).Return(u, nil)

		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		// will be populated with user in a handler
		// if AuthUser middleware is successful
		var contextUser *domain.Account

		// see this issue - https://github.com/gin-gonic/gin/issues/323
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		// we create a handler to return "user added to context" as this
		// is the only way to test modified context
		r.GET("/api/account/me", AuthUser(mockTokenService), func(c *gin.Context) {
			contextKeyVal, _ := c.Get("account")
			contextUser = contextKeyVal.(*domain.Account)
		})

		request, _ := http.NewRequest(http.MethodGet, "/api/account/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, u, contextUser)

		mockTokenService.AssertCalled(t, "ValidateAccessToken", validTokenHeader)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockTokenService.On("ValidateAccessToken", invalidTokenHeader).Return(nil, invalidTokenErr)

		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/api/account/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/api/account/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenService.AssertCalled(t, "ValidateAccessToken", invalidTokenHeader)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {

		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/api/account/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/api/account/me", http.NoBody)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenService.AssertNotCalled(t, "ValidateAccessToken")
	})
}
