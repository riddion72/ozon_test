package resolvers

import (
	"context"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/graph/model"
	"github.com/riddion72/ozon_test/internal/logger"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*domain.Post, error) {
	const f = "resolver.CreatePost"
	logger.Info("Creating post", slog.String("func", f), slog.String("title", input.Title), slog.String("user", input.User))

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
		logger.Error("Failed to create post", slog.String("func", f), slog.String("title", input.Title), slog.String("user", input.User), slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Post created successfully", slog.String("func", f), slog.Int("postID", createdPost.ID))
	return createdPost, nil
}

func (r *mutationResolver) CloseCommentsPost(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error) {
	const f = "resolver.CloseCommentsPost"
	logger.Info("Closing comments for post", slog.String("func", f), slog.Int("postID", postID), slog.String("user", user), slog.Bool("commentsAllowed", commentsAllowed))

	post, err := r.PostsSrvc.CloseComments(ctx, user, postID, commentsAllowed)
	if err != nil {
		logger.Error("Failed to close comments", slog.String("func", f), slog.Int("postID", postID), slog.String("user", user), slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Comments closed successfully", slog.String("func", f), slog.Int("postID", postID))
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*domain.Comment, error) {
	const f = "resolver.CreateComment"
	logger.Info("Creating comment", slog.String("func", f), slog.Int("postID", input.PostID), slog.String("user", input.User))

	comment := &domain.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		User:     input.User,
		Text:     input.Text,
	}

	createdComment, err := r.CommentsSrvs.Create(ctx, comment)
	if err != nil {
		logger.Error("Failed to create comment", slog.String("func", f), slog.Int("postID", input.PostID), slog.String("user", input.User), slog.String("error", err.Error()))
		return nil, err
	}

	r.Notifier.Notify(comment.PostID, createdComment)
	logger.Info("Comment created successfully", slog.String("func", f), slog.Int("commentID", createdComment.ID))
	return createdComment, nil
}
