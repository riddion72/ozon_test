package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	utils "github.com/riddion72/ozon_test/pkg/utils"
)

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	query := `INSERT INTO comments(author, post_id, parent_id, content)
	VALUES($1, $2, $3, $4)
	RETURNING *`

	err := r.db.QueryRowxContext(ctx, query,
		comment.User, comment.PostID, comment.ParentID, comment.Text,
	).Scan(&comment.ID, &comment.User, &comment.PostID, &comment.ParentID, &comment.Text, &comment.CreatedAt)

	if err != nil {
		logger.Error("Failed to create comment", slog.String("error", err.Error()))
		return nil, utils.PrepareError(err)
	}

	logger.Info("Comment created", slog.Int("commentID", comment.ID), slog.Int("postID", comment.PostID))
	return comment, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (domain.Comment, bool) {
	query := `SELECT id, author, post_id, parent_id, content, created_at 
	FROM comments 
	WHERE id = $1`
	row := r.db.QueryRowContext(context.Background(), query, id)

	var comment domain.Comment
	err := row.Scan(
		&comment.ID, &comment.User, &comment.PostID,
		&comment.ParentID, &comment.Text, &comment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Warn("Comment not found", slog.Int("commentID", id))
		return comment, false
	} else if err != nil {
		logger.Error("Error fetching comment by ID", slog.String("error", err.Error()))
		return comment, false
	}

	logger.Info("Fetched comment by ID", slog.Int("commentID", id))
	return comment, true
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error) {
	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, author, post_id, parent_id, content, created_at 
         FROM comments 
         WHERE post_id = $1
         LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)

	if err != nil {
		logger.Error("Failed to fetch comments by post ID", slog.String("error", err.Error()))
		return nil, utils.PrepareError(err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID, &comment.User, &comment.PostID, &comment.ParentID,
			&comment.Text, &comment.CreatedAt,
		); err != nil {
			logger.Error("Error scanning comment", slog.String("error", err.Error()))
			return nil, utils.PrepareError(err)
		}
		if comment.ParentID == nil {
			comments = append(comments, comment)
		}
	}
	logger.Info("Fetched comments by post ID", slog.Int("postID", postID), slog.Int("count", len(comments)))
	return comments, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, commentID int, limit int, offset int) ([]domain.Comment, error) {
	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, author, post_id, parent_id, content, created_at 
         FROM comments 
         WHERE parent_id = $1 
         LIMIT $2 OFFSET $3`,
		commentID, limit, offset,
	)

	if err != nil {
		logger.Error("Failed to fetch replies", slog.String("error", err.Error()))
		return nil, utils.PrepareError(err)
	}

	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID, &comment.User, &comment.PostID, &comment.ParentID,
			&comment.Text, &comment.CreatedAt,
		); err == nil {
			comments = append(comments, comment)
		} else {
			logger.Error("Error scanning reply", slog.String("error", err.Error()))
			return nil, utils.PrepareError(err)
		}
	}
	return comments, nil
}

func (r *CommentRepo) CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error) {
	query := `SELECT EXISTS (
        SELECT 1 FROM comments WHERE id = $1 AND post_id = $2
    )`

	var exists bool
	err := r.db.QueryRowxContext(ctx, query, commentID, postID).Scan(&exists)

	if err != nil {
		logger.Error("Failed to check if comment under post", slog.Int("postID", postID), slog.Int("commentID", commentID), slog.String("error", err.Error()))
		return false, utils.PrepareError(err)
	}

	logger.Info("Checked comment under post", slog.Int("postID", postID), slog.Int("commentID", commentID), slog.Bool("exists", exists))
	return exists, nil
}
