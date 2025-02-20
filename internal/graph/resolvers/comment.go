package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *commentResolver) User(ctx context.Context, obj *domain.Comment) (string, error) {
	return obj.User, nil
}

func (r *commentResolver) PostID(ctx context.Context, obj *domain.Comment) (string, error) {
	return obj.PostID, nil
}

func (r *commentResolver) ParentID(ctx context.Context, obj *domain.Comment) (*string, error) {
	return obj.ParentID, nil
}
