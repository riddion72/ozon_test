package service

import (
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage"
)

type CommentService struct {
	repo     storage.CommentStorage
	postRepo storage.PostStorage
	notifier *Notifier
}

func NewCommentService(
	repo storage.CommentStorage,
	postRepo storage.PostStorage,
	notifier *Notifier,
) *CommentService {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
		notifier: notifier,
	}
}

func (s *CommentService) Create(comment domain.Comment) error {
	// Проверка существования поста
	post, exists := s.postRepo.GetByID(comment.PostID)
	if !exists {
		return errors.New("post not found")
	}

	// Проверка разрешения комментариев
	if !post.CommentsAllowed {
		return errors.New("comments are disabled for this post")
	}

	// Проверка длины комментария
	if len(comment.Text) > 2000 {
		return errors.New("comment text exceeds 2000 characters")
	}

	// Проверка существования родительского комментария
	if comment.ParentID != nil {
		if _, exists := s.repo.GetByID(*comment.ParentID); !exists {
			return errors.New("parent comment not found")
		}
	}

	return s.repo.Create(comment)
}

func (s *CommentService) GetByPostID(postID string, limit, offset int) ([]domain.Comment, error) {
	return s.repo.GetByPostID(postID, limit, offset), nil
}
