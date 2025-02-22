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

func (s *commentService) Create(ctx context.Context, comment domain.Comment) error {
	// Проверка длины комментария
	if len(comment.Text) > 2000 {
		return errors.New("comment exceeds 2000 characters")
	}

	// Проверка существования поста
	post, exists := s.postRepo.GetByID(ctx, comment.PostID)
	if !exists {
		return errors.New("post not found")
	}

	// Проверка разрешения комментариев
	if !post.CommentsAllowed {
		return errors.New("comments are disabled for this post")
	}

	// Проверка родительского комментария
	if comment.ParentID != nil {
		if _, exists := s.commentRepo.GetByID(ctx, *comment.ParentID); !exists {
			return errors.New("parent comment not found")
		}
	}

	// Создание комментария
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return err
	}
	return nil
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]domain.Comment, error) {
	var dLimit *int
	var dOffset *int
	if limit == nil {
		*dLimit = defaultLimit
		limit = dLimit
	} else {
		if *limit <= 0 || *limit > maxLimit {
			*limit = defaultLimit
		}
	}
	if offset == nil {
		*dOffset = defaultOffset
		offset = dOffset
	}
	return s.commentRepo.GetByPostID(ctx, postID, *limit, *offset)
}

func (s *commentService) GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error) {
	var dLimit *int
	var dOffset *int
	if limit == nil {
		*dLimit = defaultLimit
		limit = dLimit
	} else {
		if *limit <= 0 || *limit > maxLimit {
			*limit = defaultLimit
		}
	}
	if offset == nil {
		*dOffset = defaultOffset
		offset = dOffset
	}
	return s.commentRepo.GetReplies(ctx, commentID, *limit, *offset)
}
