package resolvers

import (
	"context"

	"github.com/riddion72/ozon_test/internal/domain"
)

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *domain.Comment, error) {
	// Подписываемся на новые комментарии
	channel := r.Notifier.Subscribe(postID)

	// Возвращаем канал подписчику
	go func() {
		<-ctx.Done() // Закрываем канал, когда контекст завершен
		r.Notifier.Unsubscribe(postID, channel)
	}()

	return channel, nil
}
