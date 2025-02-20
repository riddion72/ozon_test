package resolvers

import (
	"context"
	"fmt"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *postResolver) User(ctx context.Context, obj *domain.Post) (string, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}
