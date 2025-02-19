package inmem

import (
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
)

type PostStorage struct {
	mu    sync.RWMutex
	posts map[string]*domain.Post
}

func NewPostStorage() *PostStorage {
	return &PostStorage{
		posts: make(map[string]*domain.Post),
	}
}

func (s *PostStorage) CreatePost(post *domain.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = generateID()
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	return nil
}
