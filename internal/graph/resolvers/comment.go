package resolvers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type commentResolver struct{ *Resolver }

func (r *commentResolver) User(ctx context.Context, obj *domain.Comment) (string, error) {
	logger.Info("Fetching user for comment", slog.Int("commentID", obj.ID))
	return obj.User, nil
}

func (r *commentResolver) PostID(ctx context.Context, obj *domain.Comment) (int, error) {
	logger.Info("Fetching post ID for comment", slog.Int("commentID", obj.ID))
	return obj.PostID, nil
}

func (r *commentResolver) ParentID(ctx context.Context, obj *domain.Comment) (*int, error) {
	logger.Info("Fetching parent ID for comment", slog.Int("commentID", obj.ID))
	return obj.ParentID, nil
}

func (r *commentResolver) Replies(ctx context.Context, obj *domain.Comment, limit *int, offset *int) ([]domain.Comment, error) {
	logger.Info("Fetching replies for comment", slog.Int("commentID", obj.ID), slog.String("paginete", fmt.Sprintf("limit %v offset %v", limit, offset)))
	replies, err := r.CommentsSrvs.GetReplies(ctx, obj.ID, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch replies", slog.Int("commentID", obj.ID), slog.String("error", err.Error()))
		return nil, err
	}
	return replies, nil
}
