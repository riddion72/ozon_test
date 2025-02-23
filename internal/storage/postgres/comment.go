package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	dbErrors "github.com/riddion72/ozon_test/pkg/utils/dbErorrs"
)

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepo {
	err := db.Ping()
	if err != nil {
		logger.Debug("Failed to create post",
			slog.String("error", err.Error()))
	}
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
		return nil, dbErrors.PrepareError(err)
	}

	return comment, err
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
		logger.Debug("GetByID ErrNoRows", slog.String("error", err.Error()))
		return comment, false
	} else if err != nil {
		logger.Debug("GetByID", slog.String("error", err.Error()))
		return comment, false
	}

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
		return nil, dbErrors.PrepareError(err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID, &comment.User, &comment.PostID, &comment.ParentID,
			&comment.Text, &comment.CreatedAt,
		); err != nil {
			return nil, dbErrors.PrepareError(err)
		}
		if comment.ParentID == nil {
			comments = append(comments, comment)
		}
	}
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
		return nil, dbErrors.PrepareError(err)
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
			return nil, dbErrors.PrepareError(err)
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
		return false, dbErrors.PrepareError(err)
	}

	return exists, nil
}
