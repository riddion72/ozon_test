package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage/inmemory/tools"
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

func (r *PostRepo) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	r.Lock()
	defer r.Unlock()

	post.CreatedAt = time.Now()
	post.ID = len(r.posts) + 1
	r.posts[post.ID] = *post
	return post, nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int) (domain.Post, bool) {
	r.RLock()
	defer r.RUnlock()

	post, exists := r.posts[id]
	return post, exists
}

func (r *PostRepo) List(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	r.RLock()
	defer r.RUnlock()

	end := limit + offset
	result := make([]domain.Post, 0, end)
	for i := 0; i < len(r.posts) && end > 0; i++ {
		post, exists := r.posts[i]
		if exists {
			end--
			result = append(result, post)
		}
	}
	return tools.Paginate(result, limit, offset), nil
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
