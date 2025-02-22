package service

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage"
)

type postService struct {
	repo storage.PostStorage
}

func NewPostService(repo storage.PostStorage) *postService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, post domain.Post) error {
	// Валидация данных поста
	if post.Title == "" || post.Content == "" {
		return ErrInvalidPostData
	}
	return s.repo.Create(ctx, post)
}

func (s *postService) GetByID(ctx context.Context, id string) (domain.Post, bool) {
	return s.repo.GetByID(ctx, id)
}

func (s *postService) GetPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10 // Дефолтное значение
	}
	ans := s.repo.List(ctx, limit, offset)
	//TO DO вернуть ошибку
	return ans, nil
}
