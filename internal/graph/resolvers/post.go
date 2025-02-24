package resolvers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type postResolver struct{ *Resolver }

func (r *postResolver) User(ctx context.Context, obj *domain.Post) (string, error) {
	logger.Info("Fetching user for Post", slog.Int("PostID", obj.ID))
	return obj.User, nil
}

func (r *postResolver) Comments(ctx context.Context, obj *domain.Post, limit, offset *int) ([]domain.Comment, error) {
	logger.Info("Fetching comments", slog.Int("PostID", obj.ID), slog.String("paginete", fmt.Sprintf("limit %v offset %v", limit, offset)))
	return r.CommentsSrvs.GetCommentsByPostID(ctx, obj.ID, limit, offset)
}
