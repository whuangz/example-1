package handle_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/handler"
	"github.com/whuangz/go-example/go-api/mocks"
)

func TestGetPosts(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {

		mockPost := domain.Post{
			ID:        1,
			Title:     "Test Mock 1",
			Content:   "Test Mock Desc",
			CreatedAt: time.Now(),
			Author: domain.Author{
				ID:       1,
				Username: "Whuangz",
				Email:    "whuangz@gmail.com",
			},
		}

		mockService := new(mocks.MockPostService)
		mockListPost := make([]domain.Post, 0)
		mockListPost = append(mockListPost, mockPost)

		mockService.On("FindAll", mock.AnythingOfType("*gin.Context")).Return(mockListPost, nil)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodGet, "/posts", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		respErr := domain.NewNotFound("posts", "")

		mockService := new(mocks.MockPostService)

		mockService.On("FindAll", mock.Anything).Return(nil, respErr)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodGet, "/posts", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, respErr.Status(), rec.Code)

		mockService.AssertExpectations(t)
	})
}
