package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (r *CommentRepo) transact(ctx context.Context, f func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	err = f(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	commited = true
	return nil
}
