package service

import (
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage"
)

type PostService struct {
	repo storage.PostStorage
}

func NewPostService(repo storage.PostStorage) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(post domain.Post) error {
	return s.repo.Create(post)
}

func (s *PostService) GetByID(id string) (domain.Post, bool) {
	return s.repo.GetByID(id)
}

func (s *PostService) List(limit, offset int) []domain.Post {
	return s.repo.List(limit, offset)
}
