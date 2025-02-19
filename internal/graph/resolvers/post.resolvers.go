package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input NewPost) (*domain.Post, error) {
	post := &domain.Post{
		Title:           input.Title,
		User:            input.User,
		Content:         input.Content,
		CommentsAllowed: *input.CommentsEnabled, // Обработать nullable
	}

	err := r.Storage.CreatePost(post)
	return post, err
}

func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]*domain.Post, error) {
	return r.Storage.GetPosts(limit, offset)
}
