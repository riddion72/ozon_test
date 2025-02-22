package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *domain.Comment, error) {
	// Подписываемся на новые комментарии
	channel := r.Notifier.Subscribe(postID)

	// Закрываем канал, когда контекст завершен
	go func() {
		<-ctx.Done()
		r.Notifier.Unsubscribe(postID, channel)
	}()

	return channel, nil
}
