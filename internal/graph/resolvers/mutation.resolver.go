package resolvers

import (
	"context"
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/service/"
)

type mutationResolver struct {
	services *service.Services
}

func (r *mutationResolver) CreatePost(ctx context.Context, input NewPost) (*domain.Post, error) {
	post := domain.Post{
		ID:              generateID(),
		Title:           input.Title,
		User:            input.User,
		Content:         input.Content,
		CommentsAllowed: input.CommentsEnabled,
	}

	if err := r.services.Posts.Create(post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input NewComment) (*domain.Comment, error) {
	post, exists := r.services.Posts.GetByID(input.PostId)
	if !exists {
		return nil, errors.New("post not found")
	}

	if !post.CommentsAllowed {
		return nil, errors.New("comments are disabled for this post")
	}

	if len(input.Text) > 2000 {
		return nil, errors.New("comment is too long")
	}

	comment := domain.Comment{
		ID:       generateID(),
		PostID:   input.PostId,
		ParentID: input.ParentId,
		User:     input.User,
		Text:     input.Text,
	}

	if err := r.services.Comments.Create(comment); err != nil {
		return nil, err
	}

	r.services.Notifier.Notify(comment.PostID, comment)
	return &comment, nil
}
