package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/graph/generated"
)

// Resolver корневая структура для всех резолверов
type Resolver struct {
	PostsSrvc    PostService
	CommentsSrvs CommentService
	Notifier     Notifier
}

// NewResolver конструктор для Resolver
func NewResolver(posts PostService, comments CommentService, notifier Notifier) *Resolver {
	return &Resolver{
		PostsSrvc:    posts,
		CommentsSrvs: comments,
		Notifier:     notifier,
	}
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type Notifier interface {
	Subscribe(postID int) chan *domain.Comment
	Unsubscribe(postID int, ch chan *domain.Comment)
	Notify(postID int, comment *domain.Comment)
}

// PostService интерфейс для работы с постами
type PostService interface {
	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetPostByID(ctx context.Context, id int) (domain.Post, bool)
	GetPosts(ctx context.Context, limit, offset int) ([]domain.Post, error)
	CloseComments(ctx context.Context, user string, postID int, commentsAllowed bool) (*domain.Post, error)
}

// CommentService интерфейс для работы с коментами
type CommentService interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]domain.Comment, error)
	GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]domain.Comment, error)
}
