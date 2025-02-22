package service

import (
	"context"
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage"
)

var (
	ErrCommentTooLong   = errors.New("comment exceeds 2000 characters")
	ErrCommentsDisabled = errors.New("comments are disabled for this post")
	ErrParentNotFound   = errors.New("parent comment not found")
	ErrPostNotFound     = errors.New("post not found")
	ErrInvalidPostData  = errors.New("invalid post data")
)

type commentService struct {
	commentRepo storage.CommentStorage
	postRepo    storage.PostStorage
	notifier    *Notifier
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
		return ErrCommentTooLong
	}

	// Проверка существования поста
	post, exists := s.postRepo.GetByID(ctx, comment.PostID)
	if !exists {
		return ErrPostNotFound
	}

	// Проверка разрешения комментариев
	if !post.CommentsAllowed {
		return ErrCommentsDisabled
	}

	// Проверка родительского комментария
	if comment.ParentID != nil {
		if _, exists := s.commentRepo.GetByID(ctx, *comment.ParentID); !exists {
			return ErrParentNotFound
		}
	}

	// Создание комментария
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return err
	}
	return nil
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]domain.Comment, error) {
	return s.commentRepo.GetByPostID(ctx, postID, limit, offset)
}
