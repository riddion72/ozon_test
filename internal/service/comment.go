package service

import (
	"context"
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
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
	if len(comment.Text) > 2000 {
		return nil, errors.New("comment exceeds 2000 characters")
	}

	post, exists := s.postRepo.GetByID(ctx, comment.PostID)
	if !exists {
		return nil, errors.New("post not found")
	}

	if !post.CommentsAllowed {
		return nil, errors.New("comments are disabled for this post")
	}

	if comment.ParentID != nil {
		exists, err := s.commentRepo.CheckCommentUnderPost(ctx, comment.PostID, *comment.ParentID)
		if err != nil {
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
		return nil, err
	}
	return createdComment, nil
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]domain.Comment, error) {
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
		return nil, errors.New("post not found")
	}

	return s.commentRepo.GetByPostID(ctx, postID, *limit, *offset)
}

func (s *commentService) GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error) {
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
		return []domain.Comment{}, errors.New("parent comment not found")
	}
	return s.commentRepo.GetReplies(ctx, commentID, *limit, *offset)
}
