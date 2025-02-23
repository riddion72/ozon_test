package postgres

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	dbErrors "github.com/riddion72/ozon_test/pkg/utils/dbErorrs"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	err := db.Ping()
	if err != nil {
		logger.Debug("Failed to create post",
			slog.String("error", err.Error()))
	}
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	query := `
		INSERT INTO posts (title, author, content, comments_allowed)
		VALUES ($1, $2, $3, $4)
		RETURNING *`

	err := r.db.QueryRowxContext(ctx, query,
		post.Title, post.User, post.Content, post.CommentsAllowed,
	).Scan(&post.ID, &post.Title, &post.User, &post.Content, &post.CommentsAllowed, &post.CreatedAt)

	if err != nil {
		logger.Error("Failed to create post", slog.String("error", err.Error()), slog.Int("post_id", post.ID))
		err = dbErrors.PrepareError(err)
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) GetByID(ctx context.Context, postID int) (domain.Post, bool) {
	var post domain.Post
	query := `SELECT * FROM posts WHERE id = $1`

	err := r.db.QueryRowxContext(ctx, query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.User,
		&post.Content,
		&post.CommentsAllowed,
		&post.CreatedAt)
	if err != nil {
		logger.Debug("Post not found", slog.String("error", err.Error()), slog.Int("post_id", postID))
		return domain.Post{}, false
	}
	return post, true
}

func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, title, author, content, comments_allowed, created_at 
         FROM posts LIMIT $1 OFFSET $2`,
		limit, offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(
			&post.ID, &post.Title, &post.User,
			&post.Content, &post.CommentsAllowed, &post.CreatedAt,
		); err == nil {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *PostRepository) CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error) {
	var post domain.Post
	query := `UPDATE posts SET comments_allowed = $1 WHERE id = $2 RETURNING *`

	err := r.db.QueryRowxContext(ctx, query, commentsAllowed, postID).Scan(
		&post.ID,
		&post.Title,
		&post.User,
		&post.Content,
		&post.CommentsAllowed,
		&post.CreatedAt)
	if err != nil {
		logger.Error("Failed to update comments allowed status",
			slog.String("error", err.Error()),
			slog.Int("post_id", postID))
		return nil, err
	}
	return &post, nil
}
