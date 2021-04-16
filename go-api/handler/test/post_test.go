package handle_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
		mockListPostResp := make([]domain.Post, 0)
		mockListPostResp = append(mockListPostResp, mockPost)

		mockService.On("FindAll", mock.AnythingOfType("*gin.Context")).Return(mockListPostResp, nil)

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

func TestSavePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockPostService)

	t.Run("success", func(t *testing.T) {
		mockPostResp := domain.Post{
			Title:   "Test Mock 1",
			Content: "Test Mock Desc",
		}

		j, err := json.Marshal(mockPostResp)
		assert.NoError(t, err)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("*domain.Post"),
		}
		mockService.On("Save", mockArgs...).Return(nil)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodPost, "/posts", strings.NewReader(string(j)))
		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Fail", func(t *testing.T) {
		respErr := domain.NewBadRequest("missing param")

		mockPostResp := domain.Post{
			Title: "Test Mock 1",
		}

		j, err := json.Marshal(mockPostResp)
		assert.NoError(t, err)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("*domain.Post"),
		}
		mockService.On("Save", mockArgs...).Return(respErr)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodPost, "/posts", strings.NewReader(string(j)))
		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestFindPostByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockPostService)

	t.Run("Success", func(t *testing.T) {

		mockPostResp := domain.Post{
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

		mockService.On("FindByID", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int32")).Return(mockPostResp, nil)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodGet, "/posts/1", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)

	})

	t.Run("Fail", func(t *testing.T) {
		respErr := domain.NewNotFound("post", "")

		mockService := new(mocks.MockPostService)

		mockService.On("FindByID", mock.Anything, mock.AnythingOfType("int32")).Return(nil, respErr)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodGet, "/posts/0", nil)

		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, respErr.Status(), rec.Code)

		mockService.AssertExpectations(t)
	})
}

func TestUpdatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockPostService)

	t.Run("Success", func(t *testing.T) {

		mockPostResp := domain.Post{
			Title: "Update Test Mock 1",
		}

		j, err := json.Marshal(mockPostResp)
		assert.NoError(t, err)

		mockService.On("Update", mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("int32"),
			mock.AnythingOfType("*domain.Post")).Return(nil)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodPatch, "/posts/1", strings.NewReader(string(j)))
		req.Header.Add("Content-Type", "application/json")

		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Fail", func(t *testing.T) {
		respErr := domain.NewNotFound("id", "not found")

		mockService.On("Update", mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("int32"),
			mock.AnythingOfType("*domain.Post")).Return(respErr)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodPatch, "/posts/", nil)
		req.Header.Add("Content-Type", "application/json")

		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, respErr.Status(), rec.Code)

		mockService.AssertExpectations(t)
	})
}

func TestDeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockPostService)

	t.Run("Success", func(t *testing.T) {
		mockService.On("Delete", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int32")).Return(nil)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodDelete, "/posts/1", nil)
		assert.NoError(t, err)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Fail", func(t *testing.T) {
		respErr := domain.NewNotFound("id", "not found")

		mockService.On("Delete", mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("int32")).Return(respErr)

		rec := httptest.NewRecorder()
		router := gin.New()
		handler.NewPostHandler(router, mockService)

		req, err := http.NewRequest(http.MethodDelete, "/posts/", nil)

		assert.NoError(t, err)

		router.ServeHTTP(rec, req)

		assert.Equal(t, respErr.Status(), rec.Code)

		mockService.AssertExpectations(t)
	})
}
