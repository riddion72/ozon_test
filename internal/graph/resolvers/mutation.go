package resolvers

import (
	"context"
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/graph/model"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*domain.Post, error) {
	post := domain.Post{
		Title:   input.Title,
		User:    input.User,
		Content: input.Content,
	}

	if err := r.PostsSrvc.Create(ctx, post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *mutationResolver) CloseCommentsPost(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error) {
	post, err := r.PostsSrvc.CloseComments(ctx, user, postID, commentsAllowed)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*domain.Comment, error) {
	post, exists := r.PostsSrvc.GetPostByID(ctx, input.PostID)
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
		PostID:   input.PostID,
		ParentID: input.ParentID,
		User:     input.User,
		Text:     input.Text,
	}

	if err := r.CommentsSrvs.Create(ctx, comment); err != nil {
		return nil, err
	}

	r.Notifier.Notify(comment.PostID, &comment)
	return &comment, nil
}
