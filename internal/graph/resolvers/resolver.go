package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

// Resolver корневая структура для всех резолверов
type Resolver struct {
	Posts    PostService
	Comments CommentService
	Notifier Notifier
}

// NewResolver конструктор для Resolver
func NewResolver(posts PostService,
	comments CommentService,
	notifier Notifier,
) *Resolver {
	return &Resolver{
		Posts:    posts,
		Comments: comments,
		Notifier: notifier,
	}
}

// func (r *Resolver) Complexity() graphql.ComplexityRoot {
// 	return graphql.ComplexityRoot{
// 		// ограничить глубину вложенности комментариев
// 		Comment: func(childComplexity int) int {
// 			return childComplexity * 2 // Множитель сложности для вложенных комментариев
// 		},
// 	}
// }

type Notifier interface {
	Subscribe(postID int) chan *domain.Comment
	Unsubscribe(postID int, ch chan *domain.Comment)
	Notify(postID int, comment *domain.Comment)
}

// PostService интерфейс для работы с постами
type PostService interface {
	Create(ctx context.Context, post domain.Post) error
	GetPostByID(ctx context.Context, id int) (domain.Post, bool)
	GetPosts(ctx context.Context, limit, offset int) ([]domain.Post, error)
	CloseComments(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error)
}

type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) error
	GetCommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]domain.Comment, error)
	GetReplies(commentID int, limit *int, offset *int) ([]domain.Comment, error)
}
