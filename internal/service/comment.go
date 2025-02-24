package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage"
)

const (
	defaultLimit  = 10
	maxLimit      = 100
	defaultOffset = 0
)

type commentService struct {
	commentRepo storage.CommentStorage
	postRepo    storage.PostStorage
}

func NewCommentService(
	commentRepo storage.CommentStorage,
	postRepo storage.PostStorage,
) *commentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	const f = "commentService.Create"
	logger.Info("Creating comment", slog.String("func", f), slog.Int("postID", comment.PostID))
	if len(comment.Text) > 2000 {
		return nil, errors.New("comment exceeds 2000 characters")
	}

	post, exists := s.postRepo.GetByID(ctx, comment.PostID)
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", comment.PostID))
		return nil, errors.New("post not found")
	}

	if !post.CommentsAllowed {
		return nil, errors.New("comments are disabled for this post")
	}

	if comment.ParentID != nil {
		exists, err := s.commentRepo.CheckCommentUnderPost(ctx, comment.PostID, *comment.ParentID)
		if err != nil {
			logger.Error("Error checking parent comment", slog.String("func", f), slog.Int("postID", comment.PostID), slog.String("error", err.Error()))
			return nil, err
		}
		if !exists {
			return nil, errors.New("parent comment not found")
		}
	}

	if comment.Text == "" {
		return nil, errors.New("invalid comment data")
	}

	createdComment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		logger.Error("Failed to create comment", slog.String("func", f), slog.Int("postID", comment.PostID), slog.String("error", err.Error()))
		return nil, err
	}

	return createdComment, nil
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]domain.Comment, error) {
	const f = "commentService.GetCommentsByPostID"
	logger.Info("Fetching comments by post ID", slog.String("func", f), slog.Int("postID", postID))

	defaultLimit := defaultLimit
	defaultOffset := defaultOffset
	if limit == nil {
		limit = &defaultLimit
	} else {
		if *limit < 0 || *limit > maxLimit {
			return nil, errors.New("invalid value of limit")
		}
	}

	if offset == nil {
		offset = &defaultOffset
	} else {
		if *offset < 0 {
			return nil, errors.New("invalid value of offset")
		}
	}
	_, exists := s.postRepo.GetByID(ctx, postID)
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", postID))
		return nil, errors.New("post not found")
	}
	comments, err := s.commentRepo.GetByPostID(ctx, postID, *limit, *offset)
	if err != nil {
		logger.Error("Failed to fetch comments", slog.String("func", f), slog.Int("postID", postID), slog.String("error", err.Error()))
		return nil, err
	}

	return comments, nil
}

func (s *commentService) GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error) {
	const f = "commentService.GetReplies"
	logger.Info("Fetching replies for comment", slog.String("func", f), slog.Int("commentID", commentID))

	defaultLimit := defaultLimit
	defaultOffset := defaultOffset

	if limit == nil {
		limit = &defaultLimit
	} else {
		if *limit < 0 || *limit > maxLimit {
			return nil, errors.New("invalid value of limit")
		}
	}

	if offset == nil {
		offset = &defaultOffset
	} else {
		if *offset < 0 {
			return nil, errors.New("invalid value of offset")
		}
	}
	if _, exists := s.commentRepo.GetByID(ctx, commentID); !exists {
		logger.Warn("Parent comment not found", slog.String("func", f), slog.Int("commentID", commentID))
		return []domain.Comment{}, errors.New("parent comment not found")
	}
	replies, err := s.commentRepo.GetReplies(ctx, commentID, *limit, *offset)
	if err != nil {
		logger.Error("Failed to fetch replies", slog.String("func", f), slog.Int("commentID", commentID), slog.String("error", err.Error()))
		return nil, err
	}

	return replies, nil
}
