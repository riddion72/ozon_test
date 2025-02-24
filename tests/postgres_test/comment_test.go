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

func TestCommentRepo_Create(t *testing.T) {
	logger.MustInit("prod")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewCommentRepository(sqlxDB)

	comment := &domain.Comment{
		User:     "test user",
		PostID:   1,
		ParentID: nil,
		Text:     "This is a comment",
	}

	createdAt := time.Now()
	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(comment.User, comment.PostID, comment.ParentID, comment.Text).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author", "post_id", "parent_id", "content", "created_at"}).
			AddRow(1, comment.User, comment.PostID, comment.ParentID, comment.Text, createdAt))

	createdComment, err := repo.Create(context.Background(), comment)
	require.NoError(t, err)
	assert.Equal(t, 1, createdComment.ID)
	assert.Equal(t, comment.User, createdComment.User)
	assert.Equal(t, comment.PostID, createdComment.PostID)
	assert.Equal(t, comment.Text, createdComment.Text)
}

func TestCommentRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewCommentRepository(sqlxDB)

	commentID := 1
	createdAt := time.Now()
	mock.ExpectQuery(`SELECT id, author, post_id, parent_id, content, created_at`).
		WithArgs(commentID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author", "post_id", "parent_id", "content", "created_at"}).
			AddRow(commentID, "test user", 1, nil, "This is a comment", createdAt))

	comment, exists := repo.GetByID(context.Background(), commentID)
	require.True(t, exists)
	assert.Equal(t, commentID, comment.ID)
	assert.Equal(t, "test user", comment.User)
}

func TestCommentRepo_GetByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewCommentRepository(sqlxDB)

	postID := 1
	createdAt := time.Now()
	mock.ExpectQuery(`SELECT id, author, post_id, parent_id, content, created_at`).
		WithArgs(postID, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author", "post_id", "parent_id", "content", "created_at"}).
			AddRow(1, "test user", postID, nil, "This is a comment", createdAt))

	comments, err := repo.GetByPostID(context.Background(), postID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "test user", comments[0].User)
}

func TestCommentRepo_GetReplies(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewCommentRepository(sqlxDB)

	commentID := 1
	createdAt := time.Now()
	mock.ExpectQuery(`SELECT id, author, post_id, parent_id, content, created_at`).
		WithArgs(commentID, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author", "post_id", "parent_id", "content", "created_at"}).
			AddRow(2, "reply user", 1, &commentID, "This is a reply", createdAt))

	replies, err := repo.GetReplies(context.Background(), commentID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "reply user", replies[0].User)
	assert.Equal(t, commentID, *replies[0].ParentID)
}

func TestCommentRepo_CheckCommentUnderPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := postgres.NewCommentRepository(sqlxDB)

	postID := 1
	commentID := 2
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(commentID, postID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.CheckCommentUnderPost(context.Background(), postID, commentID)
	require.NoError(t, err)
	assert.True(t, exists)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(commentID, postID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err = repo.CheckCommentUnderPost(context.Background(), postID, commentID)
	require.NoError(t, err)
	assert.False(t, exists)
}
