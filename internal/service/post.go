package service

import (
	"context"
	"errors"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage"
)

type postService struct {
	repo storage.PostStorage
}

func NewPostService(repo storage.PostStorage) *postService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	// Валидация данных поста
	if post.Title == "" || post.Content == "" {
		return nil, errors.New("invalid post data")
	}
	return s.repo.Create(ctx, post)
}

func (s *postService) GetPostByID(ctx context.Context, id int) (domain.Post, bool) {
	return s.repo.GetByID(ctx, id)
}

func (s *postService) GetPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	if limit <= 0 || limit > maxLimit {
		return []domain.Post{}, errors.New("invalid value of limit")
	}

	ans, err := s.repo.List(ctx, limit, offset)
	return ans, err
}

func (s *postService) CloseComments(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error) {
	post, exist := s.repo.GetByID(ctx, postID)
	if !exist {
		return nil, errors.New("post not found")
	}

	if post.User != user {
		return nil, errors.New("access denied")
	}

	editedPost, err := s.repo.CommentsAllowed(ctx, postID, commentsAllowed)
	if err != nil {
		return nil, err
	}

	return editedPost, nil
}
