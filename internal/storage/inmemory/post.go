package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
)

type PostRepo struct {
	sync.RWMutex
	posts map[int]domain.Post
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		posts: make(map[int]domain.Post),
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

func (r *PostRepo) GetByID(ctx context.Context, id int) (domain.Post, bool) {
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

func (r *PostRepo) CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error) {
	r.Lock()
	defer r.Unlock()

	post, exists := r.posts[postID]
	if !exists {
		return nil, errors.New("post not found")
	}

	post.CommentsAllowed = commentsAllowed
	r.posts[postID] = post
	return &post, nil
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
