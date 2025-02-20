package postgres

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post domain.Post) error {
	query := `
		INSERT INTO posts (id, title, user_id, content, comments_allowed, created_at)
		VALUES (:id, :title, :user_id, :content, :comments_allowed, :created_at)`

	_, err := r.db.NamedExecContext(context.Background(), query,
		map[string]interface{}{
			"id":               post.ID,
			"title":            post.Title,
			"user_id":          post.User,
			"content":          post.Content,
			"comments_allowed": post.CommentsAllowed,
			"created_at":       post.CreatedAt,
		})

	if err != nil {
		logger.Error("Failed to create post",
			slog.String("error", err.Error()),
			slog.String("post_id", post.ID))
		return err
	}
	return nil
}

func (r *PostRepository) GetByID(id string) (domain.Post, bool) {
	var post domain.Post
	query := `SELECT * FROM posts WHERE id = $1`

	err := r.db.Get(&post, query, id)
	if err != nil {
		logger.Debug("Post not found", slog.String("post_id", id))
		return domain.Post{}, false
	}
	return post, true
}

func (r *PostRepository) List(limit, offset int) []domain.Post {
	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, title, user_id, content, comments_allowed, created_at 
         FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)

	if err != nil {
		return nil
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
	return posts
}
