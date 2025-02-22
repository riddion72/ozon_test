package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

type commentResolver struct{ *Resolver }

func (r *commentResolver) User(ctx context.Context, obj *domain.Comment) (string, error) {
	return obj.User, nil
}

func (r *commentResolver) PostID(ctx context.Context, obj *domain.Comment) (int, error) {
	return obj.PostID, nil
}

func (r *commentResolver) ParentID(ctx context.Context, obj *domain.Comment) (*int, error) {
	return obj.ParentID, nil
}

func (r *commentResolver) Replies(ctx context.Context, obj *domain.Comment, limit *int, offset *int) ([]domain.Comment, error) {
	return r.CommentsSrvs.GetReplies(ctx, obj.ID, limit, offset)
}
