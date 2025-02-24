package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage"
)

type postService struct {
	repo storage.PostStorage
}

func NewPostService(repo storage.PostStorage) *postService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	const f = "service.CreatePost"
	logger.Info("Creating post", slog.String("func", f), slog.String("title", post.Title), slog.String("user", post.User))

	if post.Title == "" || post.Content == "" {
		logger.Warn("Invalid post data", slog.String("func", f), slog.String("title", post.Title), slog.String("content", post.Content))
		return nil, errors.New("invalid post data")
	}

	createdPost, err := s.repo.Create(ctx, post)
	if err != nil {
		logger.Error("Failed to create post", slog.String("func", f), slog.String("title", post.Title), slog.String("user", post.User), slog.String("error", err.Error()))
		return nil, err
	}
	logger.Info("Post created successfully", slog.String("func", f), slog.Int("postID", createdPost.ID))
	return createdPost, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (domain.Post, bool) {
	const f = "service.GetPostByID"
	logger.Info("Fetching post by ID", slog.String("func", f), slog.Int("postID", id))
	post, exists := s.repo.GetByID(ctx, id)
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", id))
	}
	return post, exists
}

func (s *postService) GetPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	const f = "service.GetPosts"
	logger.Info("Fetching posts", slog.String("func", f), slog.Int("limit", limit), slog.Int("offset", offset))

	if limit <= 0 || limit > maxLimit {
		logger.Warn("Invalid value of limit", slog.String("func", f), slog.Int("limit", limit))
		return []domain.Post{}, errors.New("invalid value of limit")
	}

	posts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch posts", slog.String("func", f), slog.String("error", err.Error()))
		return nil, err
	}

	return posts, nil
}

func (s *postService) CloseComments(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error) {
	const f = "service.CloseComments"
	logger.Info("Closing comments for post", slog.String("func", f), slog.Int("postID", postID), slog.String("user", user))

	post, exist := s.repo.GetByID(ctx, postID)
	if !exist {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", postID))
		return nil, errors.New("post not found")
	}

	if post.User != user {
		logger.Warn("Access denied", slog.String("func", f), slog.Int("postID", postID), slog.String("user", user))
		return nil, errors.New("access denied")
	}

	editedPost, err := s.repo.CommentsAllowed(ctx, postID, commentsAllowed)
	if err != nil {
		logger.Error("Failed to update comments allowed", slog.String("func", f), slog.Int("postID", postID), slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Comments updated successfully", slog.String("func", f), slog.Int("postID", postID), slog.Bool("commentsAllowed", commentsAllowed))
	return editedPost, nil
}
