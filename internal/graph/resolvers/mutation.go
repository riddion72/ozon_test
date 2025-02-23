package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/graph/model"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*domain.Post, error) {
	post := &domain.Post{
		Title:   input.Title,
		User:    input.User,
		Content: input.Content,
	}
	if input.CommentsAllowed != nil {
		post.CommentsAllowed = *input.CommentsAllowed
	}

	createdPost, err := r.PostsSrvc.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	return createdPost, nil
}

func (r *mutationResolver) CloseCommentsPost(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error) {
	post, err := r.PostsSrvc.CloseComments(ctx, user, postID, commentsAllowed)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*domain.Comment, error) {
	comment := &domain.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		User:     input.User,
		Text:     input.Text,
	}

	createdComment, err := r.CommentsSrvs.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	r.Notifier.Notify(comment.PostID, createdComment)
	return createdComment, nil
}
