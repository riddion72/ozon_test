package resolvers

import (
	"context"
	"fmt"

	"github.com/riddion72/ozon_test/internal/domain"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]domain.Post, error) {
	return r.PostsSrvc.GetPosts(ctx, limit, offset)
}

func (r *queryResolver) Post(ctx context.Context, id int) (*domain.Post, error) {
	post, exists := r.PostsSrvc.GetPostByID(ctx, id)
	if !exists {
		return nil, fmt.Errorf("post not found")
	}
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID int, limit *int, offset *int) ([]domain.Comment, error) {
	return r.CommentsSrvs.GetCommentsByPostID(ctx, postID, limit, offset)
}

func (r *queryResolver) Replies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error) {
	return r.CommentsSrvs.GetReplies(ctx, commentID, limit, offset)
}
