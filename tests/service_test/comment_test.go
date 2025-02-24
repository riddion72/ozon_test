package service_test

import (
	"context"
	"testing"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCommentStorage struct {
	mock.Mock
}

func (m *MockCommentStorage) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentStorage) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func (m *MockCommentStorage) GetByID(ctx context.Context, commentID int) (domain.Comment, bool) {
	args := m.Called(ctx, commentID)
	return args.Get(0).(domain.Comment), args.Bool(1)
}

func (m *MockCommentStorage) CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error) {
	args := m.Called(ctx, postID, commentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCommentStorage) GetReplies(ctx context.Context, commentID int, limit, offset int) ([]domain.Comment, error) {
	args := m.Called(ctx, commentID, limit, offset)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func ptr(i int) *int { return &i }

func TestCommentService(t *testing.T) {
	mockCommentRepo := new(MockCommentStorage)
	mockPostRepo := new(MockPostStorage)
	commentService := service.NewCommentService(mockCommentRepo, mockPostRepo)
	ctx := context.Background()

	t.Run("Create comment", func(t *testing.T) {
		post := domain.Post{ID: 1, CommentsAllowed: true}
		comment := &domain.Comment{
			PostID:   post.ID,
			Text:     "This is a comment",
			ParentID: nil,
		}

		mockPostRepo.On("GetByID", mock.Anything, post.ID).Return(post, true)
		mockCommentRepo.On("Create", mock.Anything, comment).Return(comment, nil)

		createdComment, err := commentService.Create(ctx, comment)

		assert.NoError(t, err)
		assert.NotNil(t, createdComment)
		assert.Equal(t, comment.Text, createdComment.Text)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Create comment for !exist post", func(t *testing.T) {
		comment := &domain.Comment{
			PostID:   123,
			Text:     "This is a comment",
			ParentID: nil,
		}

		mockPostRepo.On("GetByID", mock.Anything, comment.PostID).Return(domain.Post{}, false)

		createdComment, err := commentService.Create(ctx, comment)

		assert.Error(t, err)
		assert.Nil(t, createdComment)
		assert.Equal(t, "post not found", err.Error())
		mockPostRepo.AssertExpectations(t)
	})

	t.Run("Create comment invalid data", func(t *testing.T) {
		comment := &domain.Comment{
			PostID:   1,
			Text:     "",
			ParentID: nil,
		}

		createdComment, err := commentService.Create(ctx, comment)

		assert.Error(t, err)
		assert.Nil(t, createdComment)
		assert.Equal(t, "invalid comment data", err.Error())
	})

	t.Run("Get comments by post ID", func(t *testing.T) {
		post := domain.Post{ID: 1}
		comments := []domain.Comment{
			{ID: 1, PostID: post.ID, Text: "Comment 1"},
			{ID: 2, PostID: post.ID, Text: "Comment 2"},
		}
		limit := 10
		offset := 0

		mockPostRepo.On("GetByID", mock.Anything, post.ID).Return(post, true)
		mockCommentRepo.On("GetByPostID", mock.Anything, post.ID, limit, offset).Return(comments, nil)

		retrievedComments, err := commentService.GetCommentsByPostID(ctx, post.ID, &limit, &offset)

		assert.NoError(t, err)
		assert.Len(t, retrievedComments, 2)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Get comments !exist post", func(t *testing.T) {
		postID := 123
		limit := 10
		offset := 0

		mockPostRepo.On("GetByID", mock.Anything, postID).Return(domain.Post{}, false)

		retrievedComments, err := commentService.GetCommentsByPostID(ctx, postID, &limit, &offset)

		assert.Error(t, err)
		assert.Nil(t, retrievedComments)
		assert.Equal(t, "post not found", err.Error())
	})

	t.Run("Successfully get replies", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		mockCommentRepo.On("GetByID", ctx, 1).Return(domain.Comment{ID: 1}, true)

		replies := []domain.Comment{
			{ID: 2, ParentID: ptr(1)},
			{ID: 3, ParentID: ptr(1)},
		}
		mockCommentRepo.On("GetReplies", ctx, 1, 10, 0).Return(replies, nil)

		result, err := service.GetReplies(ctx, 1, nil, nil)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		mockCommentRepo.AssertExpectations(t)
	})

	t.Run("Empty replies list", func(t *testing.T) {
		service1 := service.NewCommentService(mockCommentRepo, mockPostRepo)

		mockCommentRepo.On("GetByID", ctx, 1).Return(domain.Comment{ID: 1}, true)
		mockCommentRepo.On("GetReplies", ctx, 1, 10, 0).Return([]domain.Comment{}, nil)

		_, err := service1.GetReplies(ctx, 1, nil, nil)

		require.NoError(t, err)
	})

	t.Run("Parent comment not found", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		mockCommentRepo.On("GetByID", ctx, 123).Return(domain.Comment{}, false)

		_, err := service.GetReplies(ctx, 123, nil, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "parent comment not found")
	})

	t.Run("Invalid limit parameter", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		invalidLimit := -5
		_, err := service.GetReplies(ctx, 1, &invalidLimit, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "invalid value of limit")
	})

	t.Run("Pagination with custom parameters", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		limit := 5
		offset := 2

		mockCommentRepo.On("GetByID", ctx, 1).Return(domain.Comment{ID: 1}, true)
		mockCommentRepo.On("GetReplies", ctx, 1, 5, 2).Return([]domain.Comment{}, nil)

		_, err := service.GetReplies(ctx, 1, &limit, &offset)

		require.NoError(t, err)
		mockCommentRepo.AssertCalled(t, "GetReplies", ctx, 1, 5, 2)
	})

	t.Run("Exceed max limit", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		overLimit := 150
		_, err := service.GetReplies(ctx, 1, &overLimit, nil)

		require.Error(t, err)
		assert.EqualError(t, err, "invalid value of limit")
	})

	t.Run("Default parameters when nil", func(t *testing.T) {
		service := service.NewCommentService(mockCommentRepo, mockPostRepo)

		mockCommentRepo.On("GetByID", ctx, 1).Return(domain.Comment{ID: 1}, true)
		mockCommentRepo.On("GetReplies", ctx, 1, 10, 0).Return([]domain.Comment{}, nil)

		_, err := service.GetReplies(ctx, 1, nil, nil)

		require.NoError(t, err)
		mockCommentRepo.AssertCalled(t, "GetReplies", ctx, 1, 10, 0)
	})

}
