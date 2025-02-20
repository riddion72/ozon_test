package resolvers

import (
	"context"
	"fmt"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]domain.Post, error) {
	return r.services.Posts.List(ctx, limit, offset)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*domain.Post, error) {
	post, exists := r.services.Posts.GetByID(ctx, id)
	if !exists {
		return nil, fmt.Errorf("post not found")
	}
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID string, limit int, offset int) ([]domain.Comment, error) {
	return r.services.Comments.GetByPostID(ctx, postID, limit, offset)
}
