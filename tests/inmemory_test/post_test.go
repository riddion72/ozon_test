package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	repo := inmemory.NewPostRepo()
	ctx := context.Background()

	t.Run("creation", func(t *testing.T) {
		post := &domain.Post{
			Title:           "Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		createdPost, err := repo.Create(ctx, post)

		assert.NoError(t, err)
		assert.NotNil(t, createdPost)
		assert.Equal(t, post.Title, createdPost.Title)
		assert.Equal(t, post.User, createdPost.User)
		assert.Equal(t, post.Content, createdPost.Content)
		assert.True(t, createdPost.CommentsAllowed)
		assert.True(t, createdPost.ID == 1)
		assert.WithinDuration(t, time.Now(), createdPost.CreatedAt, time.Second)
	})

	t.Run("2 post creation", func(t *testing.T) {
		post := &domain.Post{
			Title:           "Duplicate Title",
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		}

		_, err := repo.Create(ctx, post)
		assert.NoError(t, err)

		created, err := repo.Create(context.Background(), post)
		assert.Equal(t, 3, created.ID)

	})
}

func TestGetPostByID(t *testing.T) {
	repo := inmemory.NewPostRepo()
	ctx := context.Background()

	post := &domain.Post{
		Title:           "Title",
		User:            "Author",
		Content:         "Content",
		CommentsAllowed: true,
	}

	createdPost, _ := repo.Create(ctx, post)

	t.Run("existing post", func(t *testing.T) {
		retrievedPost, exists := repo.GetByID(ctx, createdPost.ID)

		assert.True(t, exists, "Post exist")
		assert.Equal(t, createdPost, &retrievedPost, "Retrieved post should match the created post")
	})

	t.Run("non-existing post", func(t *testing.T) {
		retrievedPost, exists := repo.GetByID(ctx, 123)

		assert.False(t, exists, "Post !exist")
		assert.Equal(t, domain.Post{}, retrievedPost, "Retrieved post should be empty")
	})
}

func TestList(t *testing.T) {
	repo := inmemory.NewPostRepo()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = repo.Create(ctx, &domain.Post{
			Title:           "Title " + string(rune('A'+i-1)),
			User:            "Author",
			Content:         "Content",
			CommentsAllowed: true,
		})
	}

	t.Run("list paginate1", func(t *testing.T) {
		posts, err := repo.List(ctx, 3, 0)

		assert.NoError(t, err, "Should not return an error")
		assert.Len(t, posts, 3)
	})

	t.Run("list paginate2", func(t *testing.T) {
		posts, err := repo.List(ctx, 10, 0)

		assert.NoError(t, err, "Should not return an error")
		assert.Len(t, posts, 4)
	})

	t.Run("list offset", func(t *testing.T) {
		posts, err := repo.List(ctx, 2, 2)

		assert.NoError(t, err, "Should not return an error")
		assert.Len(t, posts, 2)
		assert.Equal(t, "Title B", posts[0].Title)
		assert.Equal(t, "Title C", posts[1].Title)
	})
}

func TestCommentsAllowed(t *testing.T) {
	repo := inmemory.NewPostRepo()
	ctx := context.Background()

	post := &domain.Post{
		Title:           "Title",
		User:            "Author",
		Content:         "Content",
		CommentsAllowed: true,
	}

	createdPost, _ := repo.Create(ctx, post)

	t.Run("update to false", func(t *testing.T) {
		updatedPost, err := repo.CommentsAllowed(ctx, createdPost.ID, false)

		assert.NoError(t, err, "Should not return an error")
		assert.False(t, updatedPost.CommentsAllowed)
	})

	t.Run("update to true", func(t *testing.T) {
		updatedPost, err := repo.CommentsAllowed(ctx, createdPost.ID, true)

		assert.NoError(t, err, "Should not return an error")
		assert.True(t, updatedPost.CommentsAllowed)
	})

	t.Run("update !exist post", func(t *testing.T) {
		_, err := repo.CommentsAllowed(ctx, 123, false)

		assert.Error(t, err, "Should return an error for non-existing post")
		assert.Equal(t, "post not found", err.Error())
	})
}
