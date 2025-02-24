package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage/postgres"
	"github.com/riddion72/ozon_test/pkg/testhelper"
)

func TestCreatePost(t *testing.T) {
	db := testhelper.SetupDBMock(t)

	repo := postgres.NewPostRepository(db)

	post := &domain.Post{
		Title:           "Test Title",
		User:            "Test Author",
		Content:         "Test Content",
		CommentsAllowed: true,
	}

	// Ожидаем, что будет выполнен запрос на вставку поста
	// query := `
	// 	INSERT INTO posts (title, author, content, comments_allowed)
	// 	VALUES ($1, $2, $3, $4)
	// 	RETURNING *`

	// Настраиваем ожидания для мока
	// dbMock := sqlxmock.Newx()
	// dbMock.ExpectQuery(query).
	// 	WithArgs(post.Title, post.User, post.Content, post.CommentsAllowed).
	// 	WillReturnRows(sqlxmock.NewRows([]string{"id", "title", "author", "content", "comments_allowed", "created_at"}).
	// 		AddRow(1, post.Title, post.User, post.Content, post.CommentsAllowed, "2023-01-01 00:00:00"))

	// Выполняем создание поста
	createdPost, err := repo.Create(context.Background(), post)

	// Проверяем, что ошибок нет
	require.NoError(t, err)
	// Проверяем, что пост был создан с правильными данными
	assert.Equal(t, post.Title, createdPost.Title)
	assert.Equal(t, post.User, createdPost.User)

	// Проверяем, что все ожидания были выполнены
	// assert.NoError(t, dbMock.ExpectationsWereMet())
}
