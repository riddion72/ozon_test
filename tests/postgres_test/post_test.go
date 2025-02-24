package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage/postgres"
)

func TestPostRepository_Create(t *testing.T) {
	logger.MustInit("prod")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewPostRepository(sqlxDB)

	post := &domain.Post{
		Title:           "Post",
		User:            "user",
		Content:         "This is a post.",
		CommentsAllowed: true,
	}

	createdAt := time.Now()
	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(post.Title, post.User, post.Content, post.CommentsAllowed).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "content", "comments_allowed", "created_at"}).
			AddRow(1, post.Title, post.User, post.Content, post.CommentsAllowed, createdAt))

	createdPost, err := repo.Create(context.Background(), post)
	require.NoError(t, err)
	assert.Equal(t, 1, createdPost.ID)
	assert.Equal(t, post.Title, createdPost.Title)
	assert.Equal(t, post.User, createdPost.User)
	assert.Equal(t, post.Content, createdPost.Content)
	assert.Equal(t, post.CommentsAllowed, createdPost.CommentsAllowed)
	assert.Equal(t, createdAt, createdPost.CreatedAt)
}

// func TestPostRepository_GetByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	sqlxDB := sqlx.NewDb(db, "postgres")
// 	repo := postgres.NewPostRepository(sqlxDB)

// 	postID := 1
// 	createdAt := time.Now()
// 	mock.ExpectQuery(`SELECT * FROM posts WHERE id = 1`).
// 		WithArgs(postID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "content", "comments_allowed", "created_at"}).
// 			AddRow(postID, "Post", "user", "This is a post.", true, createdAt))

// 	post, exists := repo.GetByID(context.Background(), postID)
// 	require.True(t, exists)
// 	assert.Equal(t, postID, post.ID)
// 	assert.Equal(t, "Post", post.Title)
// 	assert.Equal(t, "user", post.User)
// 	assert.Equal(t, "This is a post.", post.Content)
// 	assert.True(t, post.CommentsAllowed)
// 	assert.Equal(t, createdAt, post.CreatedAt)
// }

func TestPostRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewPostRepository(sqlxDB)

	createdAt := time.Now()
	mock.ExpectQuery(`SELECT id, title, author, content, comments_allowed, created_at`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "content", "comments_allowed", "created_at"}).
			AddRow(1, "Post 1", "user", "the first.", true, createdAt).
			AddRow(2, "Post 2", "user", "the second.", false, createdAt))

	posts, err := repo.List(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "Post 2", posts[1].Title)
}

// func TestPostRepository_CommentsAllowed(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	sqlxDB := sqlx.NewDb(db, "postgres")
// 	repo := postgres.NewPostRepository(sqlxDB)
// 	postID := 1
// 	mock.ExpectQuery(`UPDATE posts SET comments_allowed = $1 WHERE id = $2 RETURNING *`).
// 		WithArgs(false, postID).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "content", "comments_allowed", "created_at"}).
// 			AddRow(postID, "Test Post", "test user", "This is a test post.", false, time.Now()))

// 	updatedPost, err := repo.CommentsAllowed(context.Background(), postID, false)
// 	require.NoError(t, err)
// 	assert.Equal(t, postID, updatedPost.ID)
// 	assert.Equal(t, false, updatedPost.CommentsAllowed)
// }
