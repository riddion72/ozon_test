package storage

import "github.com/riddion72/ozon_test/internal/domain"

type PostStorage interface {
	CreatePost(post *domain.Post) error
	GetPosts(limit, offset int) ([]*domain.Post, error)
	GetPostByID(id string) (*domain.Post, error)
}

type CommentStorage interface {
	CreateComment(comment *domain.Comment) error
	GetCommentsByPost(postID string, limit, offset int) ([]*domain.Comment, error)
}
