package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/mocks"
	"github.com/whuangz/go-example/go-api/service"
)

func TestSave(t *testing.T) {
	mockRepo := new(mocks.MockPostRepo)
	mockPost := domain.Post{
		Title:   "Test Mock 1",
		Content: "Test Mock Desc",
	}

	t.Run("Success", func(t *testing.T) {

		tempMockPost := mockPost

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&tempMockPost,
		}
		mockRepo.On("Save", mockArgs...).Return(nil).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Save(ctx, &tempMockPost)
		assert.NoError(t, err)

		assert.Equal(t, tempMockPost.AuthorID, int32(1))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing Param", func(t *testing.T) {

		tempMockPost := mockPost
		tempMockPost.Title = ""

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&tempMockPost,
		}
		mockRepo.On("Save", mockArgs...).Return(domain.NewBadRequest("missing title")).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Save(ctx, &tempMockPost)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestFindAll(t *testing.T) {
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

		mockListPostResp := make([]domain.Post, 0)
		mockListPostResp = append(mockListPostResp, mockPost)

		mockRepository := new(mocks.MockPostRepo)
		mockRepository.On("FindAll", mock.Anything).Return(mockListPostResp, nil).Once()

		ps := service.NewPostService(mockRepository)

		ctx := context.TODO()
		u, err := ps.FindAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, u, mockListPostResp)
		mockRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepository := new(mocks.MockPostRepo)
		mockRepository.On("FindAll", mock.Anything).Return(nil, fmt.Errorf("Some error down the call chain")).Once()

		ps := service.NewPostService(mockRepository)

		ctx := context.TODO()
		u, err := ps.FindAll(ctx)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockRepository.AssertExpectations(t)
	})
}

func TestFindById(t *testing.T) {

	mockRepository := new(mocks.MockPostRepo)
	mockPost := domain.Post{
		// ID:        1,
		Title:     "Test Mock 1",
		Content:   "Test Mock Desc",
		CreatedAt: time.Now(),
		Author: domain.Author{
			ID:       1,
			Username: "Whuangz",
			Email:    "whuangz@gmail.com",
		},
	}

	t.Run("Success", func(t *testing.T) {

		mockRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int32")).Return(mockPost, nil).Once()

		ps := service.NewPostService(mockRepository)

		ctx := context.TODO()
		p, err := ps.FindByID(ctx, mockPost.ID)

		assert.NoError(t, err)
		assert.NotNil(t, p)

		mockRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepository.On("FindByID", mock.Anything, mock.AnythingOfType("int32")).Return(domain.Post{}, domain.NewNotFound("id", "id")).Once()

		ps := service.NewPostService(mockRepository)

		ctx := context.TODO()
		p, err := ps.FindByID(ctx, mockPost.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Post{}, p)
		mockRepository.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(mocks.MockPostRepo)
	var id int32 = 1
	mockUpdatePostParam := domain.Post{
		Title:   "Update Test Mock 2",
		Content: "Test Mock Desc",
	}

	t.Run("Success", func(t *testing.T) {

		tempMockPost := mockUpdatePostParam

		mockUpdatePostResp := domain.Post{
			ID:      1,
			Title:   "Update Test Mock 2",
			Content: "Test Mock Desc",
		}

		mockArgs := mock.Arguments{
			mock.Anything,
			mock.AnythingOfType("int32"),
			mock.AnythingOfType("*domain.Post"),
		}

		mockRepo.On("Update", mockArgs...).Return(nil).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Update(ctx, id, &tempMockPost)

		assert.NoError(t, err)
		assert.Equal(t, tempMockPost.Title, mockUpdatePostResp.Title)
		assert.NotNil(t, tempMockPost.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		errResp := domain.NewNotFound("id", "not found")
		tempMockPost := mockUpdatePostParam

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("int32"),
			mock.AnythingOfType("*domain.Post"),
		}
		mockRepo.On("Update", mockArgs...).Return(errResp).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Update(ctx, 1, &tempMockPost)

		assert.Error(t, err)
		assert.Equal(t, domain.Status(err), domain.Status(errResp))
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(mocks.MockPostRepo)

	t.Run("Success", func(t *testing.T) {

		var id int32 = 1

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("int32"),
		}

		mockRepo.On("Delete", mockArgs...).Return(nil).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("int32"),
		}
		mockRepo.On("Delete", mockArgs...).Return(domain.NewNotFound("id", "not found")).Once()

		ps := service.NewPostService(mockRepo)

		ctx := context.TODO()
		err := ps.Delete(ctx, 1)
		print(err)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
