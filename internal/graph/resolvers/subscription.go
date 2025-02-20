package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *domain.Comment, error) {
	ch := make(chan *domain.Comment, 1)

	unsubscribe := r.services.Notifier.Subscribe(postID, func(comment domain.Comment) {
		ch <- &comment
	})

	go func() {
		<-ctx.Done()
		unsubscribe()
	}()

	return ch, nil
}
