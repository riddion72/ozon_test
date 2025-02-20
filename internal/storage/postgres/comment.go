package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/domain"
)

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(comment domain.Comment) error {
	_, err := r.db.ExecContext(
		context.Background(),
		`INSERT INTO comments(id, post_id, parent_id, user_id, text, created_at)
         VALUES($1, $2, $3, $4, $5, $6)`,
		comment.ID, comment.PostID, comment.ParentID, comment.User,
		comment.Text, comment.CreatedAt,
	)
	return err
}

func (r *CommentRepo) GetByPostID(postID string, limit, offset int) []domain.Comment {
	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, post_id, parent_id, user_id, text, created_at 
         FROM comments 
         WHERE post_id = $1 
         ORDER BY created_at DESC 
         LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)

	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.ParentID,
			&comment.User, &comment.Text, &comment.CreatedAt,
		); err == nil {
			comments = append(comments, comment)
		}
	}
	return comments
}

func (r *CommentRepo) GetByID(id string) (domain.Comment, bool) {
	row := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, post_id, parent_id, user_id, text, created_at 
         FROM comments 
         WHERE id = $1`,
		id,
	)

	var comment domain.Comment
	err := row.Scan(
		&comment.ID, &comment.PostID, &comment.ParentID,
		&comment.User, &comment.Text, &comment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return comment, false
	} else if err != nil {
		return comment, false
	}

	return comment, true
}
