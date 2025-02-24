package inmemory

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
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
	const f = "inmemory.PostRepo.Create"
	r.Lock()
	defer r.Unlock()

	post.CreatedAt = time.Now()
	post.ID = len(r.posts) + 1
	r.posts[post.ID] = *post

	logger.Info("Post created", slog.String("func", f), slog.Int("postID", post.ID), slog.String("title", post.Title))
	return post, nil
}

func (r *PostRepo) GetByID(ctx context.Context, id int) (domain.Post, bool) {
	const f = "inmemory.PostRepo.GetByID"
	r.RLock()
	defer r.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", id)) // Логирование, если пост не найден
	}
	return post, exists
}

func (r *PostRepo) List(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	const f = "inmemory.PostRepo.List"
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
	logger.Info("Fetched posts", slog.String("func", f), slog.Int("count", len(result)))
	return tools.Paginate(result, limit, offset), nil
}

func (r *PostRepo) CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error) {
	const f = "inmemory.PostRepo.CommentsAllowed"
	r.Lock()
	defer r.Unlock()

	post, exists := r.posts[postID]
	if !exists {
		logger.Warn("Post not found", slog.String("func", f), slog.Int("postID", postID))
		return nil, errors.New("post not found")
	}

	post.CommentsAllowed = commentsAllowed
	r.posts[postID] = post

	logger.Info("Updated comments allowed", slog.String("func", f), slog.Int("postID", postID), slog.Bool("commentsAllowed", commentsAllowed))
	return &post, nil
}
