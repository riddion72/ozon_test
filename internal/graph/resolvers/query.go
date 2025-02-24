package resolvers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]domain.Post, error) {
	const f = "resolvers.Post"
	logger.Info("Fetching posts", slog.String("func", f), slog.Int("limit", limit), slog.Int("offset", offset))
	posts, err := r.PostsSrvc.GetPosts(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch posts", slog.String("error", err.Error()))
		return nil, err
	}
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id int) (*domain.Post, error) {
	const f = "resolvers.Post"
	logger.Info("Fetching post by ID", slog.String("func", f), slog.Int("postID", id))
	post, exists := r.PostsSrvc.GetPostByID(ctx, id)
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", id))
		return nil, fmt.Errorf("post not found")
	}
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID int, limit *int, offset *int) ([]domain.Comment, error) {
	const f = "resolvers.Comments"
	logger.Info("Fetching comments for post", slog.String("func", f), slog.Int("postID", postID),
		slog.String("paginete", fmt.Sprintf("limit %v offset %v", limit, offset)))
	comments, err := r.CommentsSrvs.GetCommentsByPostID(ctx, postID, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch comments", slog.String("func", f), slog.Int("postID", postID),
			slog.String("error", err.Error()))
		return nil, err
	}
	return comments, nil
}

func (r *queryResolver) Replies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error) {
	const f = "resolvers.Replies"
	logger.Info("Fetching replies for comment", slog.String("func", f), slog.Int("commentID", commentID),
		slog.String("paginete", fmt.Sprintf("limit %v offset %v", limit, offset)))
	replies, err := r.CommentsSrvs.GetReplies(ctx, commentID, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch replies", slog.String("func", f), slog.Int("commentID", commentID),
			slog.String("error", err.Error()))
		return nil, err
	}
	return replies, nil
}
