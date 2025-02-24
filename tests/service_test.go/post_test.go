package service_test

import (
	"context"
	"testing"

	"github.com/riddion72/ozon_test/internal/domain" // Импортируйте ваш репозиторий
	"github.com/riddion72/ozon_test/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostStorage struct {
	mock.Mock
}

func (m *MockPostStorage) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostStorage) GetByID(ctx context.Context, id int) (domain.Post, bool) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Post), args.Bool(1)
}

func (m *MockPostStorage) List(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Post), args.Error(1)
}

func (m *MockPostStorage) CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error) {
	args := m.Called(ctx, postID, commentsAllowed)
	return args.Get(0).(*domain.Post), args.Error(1)
}

func TestPostService(t *testing.T) {
	mockRepo := new(MockPostStorage)
	service := service.NewPostService(mockRepo)
	ctx := context.Background()

	t.Run("Create post", func(t *testing.T) {
		post := &domain.Post{
			Title:           "Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		mockRepo.On("Create", mock.Anything, post).Return(post, nil)

		createdPost, err := service.Create(ctx, post)

		assert.NoError(t, err)
		assert.NotNil(t, createdPost)
		assert.Equal(t, post.Title, createdPost.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create invalid post1", func(t *testing.T) {
		post := &domain.Post{
			Title:           "",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		createdPost, err := service.Create(ctx, post)

		assert.Error(t, err)
		assert.Nil(t, createdPost)
		assert.Equal(t, "invalid post data", err.Error())
	})

	t.Run("Create invalid post2", func(t *testing.T) {
		post := &domain.Post{
			Title:           "",
			User:            "",
			Content:         "",
			CommentsAllowed: true,
		}

		createdPost, err := service.Create(ctx, post)

		assert.Error(t, err)
		assert.Nil(t, createdPost)
		assert.Equal(t, "invalid post data", err.Error())
	})

	t.Run("GetPostByID exist post", func(t *testing.T) {
		post := domain.Post{
			ID:              1,
			Title:           "Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		mockRepo.On("GetByID", mock.Anything, post.ID).Return(post, true)

		retrievedPost, exists := service.GetPostByID(ctx, post.ID)

		assert.True(t, exists)
		assert.Equal(t, post, retrievedPost)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPostByID !exist post", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, 123).Return(domain.Post{}, false)

		retrievedPost, exists := service.GetPostByID(ctx, 123)

		assert.False(t, exists)
		assert.Equal(t, domain.Post{}, retrievedPost)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPosts limit1", func(t *testing.T) {
		posts := []domain.Post{
			{ID: 1, Title: "Post 1"},
			{ID: 2, Title: "Post 2"},
		}

		mockRepo.On("List", mock.Anything, 2, 0).Return(posts, nil)

		retrievedPosts, err := service.GetPosts(ctx, 2, 0)

		assert.NoError(t, err)
		assert.Len(t, retrievedPosts, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPosts limit2", func(t *testing.T) {
		posts := []domain.Post{
			{ID: 1, Title: "Post 1"},
			{ID: 2, Title: "Post 2"},
		}

		mockRepo.On("List", mock.Anything, 2, 0).Return(posts, nil)

		retrievedPosts, err := service.GetPosts(ctx, -1, 0)

		assert.Error(t, err)
		assert.Equal(t, "invalid value of limit", err.Error())
		assert.Len(t, retrievedPosts, 0)
	})

	t.Run("Close exist post", func(t *testing.T) {
		post := domain.Post{
			ID:              1,
			Title:           "Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: false,
		}

		mockRepo.On("GetByID", mock.Anything, post.ID).Return(post, true)
		mockRepo.On("CommentsAllowed", mock.Anything, post.ID, false).Return(&post, nil)

		editedPost, err := service.CloseComments(ctx, post.User, post.ID, false)

		assert.NoError(t, err)
		assert.False(t, editedPost.CommentsAllowed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("access denied", func(t *testing.T) {
		post := domain.Post{
			ID:              1,
			Title:           "Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		mockRepo.On("GetByID", mock.Anything, post.ID).Return(post, true)

		editedPost, err := service.CloseComments(ctx, "Another", post.ID, false)

		assert.Error(t, err)
		assert.Nil(t, editedPost)
		assert.Equal(t, "access denied", err.Error())
	})

	t.Run("CloseComments !exist post", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, 123).Return(domain.Post{}, false)

		editedPost, err := service.CloseComments(ctx, "Author", 123, false)

		assert.Error(t, err)
		assert.Nil(t, editedPost)
		assert.Equal(t, "post not found", err.Error())
	})
}
