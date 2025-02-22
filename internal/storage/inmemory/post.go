package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
)

type PostRepo struct {
	sync.RWMutex
	posts map[uint]domain.Post
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		posts: make(map[uint]domain.Post),
	}
}

func (r *PostRepo) Create(ctx context.Context, post domain.Post) error {
	r.Lock()
	defer r.Unlock()

	post.CreatedAt = time.Now()
	post.ID = len(r.posts)
	r.posts[post.ID] = post
	return nil
}

func (r *PostRepo) GetByID(ctx context.Context, id string) (domain.Post, bool) {
	r.RLock()
	defer r.RUnlock()

	post, exists := r.posts[id]
	return post, exists
}

func (r *PostRepo) List(ctx context.Context, limit, offset int) []domain.Post {
	r.RLock()
	defer r.RUnlock()

	result := make([]domain.Post, 0, limit)
	for _, post := range r.posts {
		result = append(result, post)
	}
	return paginate(result, limit, offset)
}

func paginate[T any](slice []T, limit, offset int) []T {
	if offset > len(slice) {
		return []T{}
	}
	end := offset + limit
	if end > len(slice) {
		end = len(slice)
	}
	return slice[offset:end]
}
