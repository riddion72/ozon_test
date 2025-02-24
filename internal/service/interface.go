package service

// import (
// 	"context"

// 	"github.com/riddion72/ozon_test/internal/domain"
// )

// type PostStorage interface {
// 	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
// 	GetByID(ctx context.Context, postID int) (domain.Post, bool)
// 	List(ctx context.Context, limit, offset int) ([]domain.Post, error)
// 	CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error)
// }

// type CommentStorage interface {
// 	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
// 	GetByID(ctx context.Context, commentID int) (domain.Comment, bool)
// 	GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error)
// 	GetReplies(ctx context.Context, commentID int, limit int, offset int) ([]domain.Comment, error)
// 	CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error)
// }
