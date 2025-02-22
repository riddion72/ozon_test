package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

type postResolver struct{ *Resolver }

func (r *postResolver) User(ctx context.Context, obj *domain.Post) (string, error) {
	return obj.User, nil
}

func (r *postResolver) Comments(ctx context.Context, obj *domain.Post, limit, offset *int) ([]domain.Comment, error) {
	return r.CommentsSrvs.GetCommentsByPostID(ctx, obj.ID, limit, offset)
}
