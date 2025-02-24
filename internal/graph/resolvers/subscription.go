package resolvers

import (
	"context"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *domain.Comment, error) {
	const f = "resolvers.Subscription.CommentAdded"
	logger.Info("Subscribing to comments for post", slog.String("func ", f), slog.Int("postID", postID))
	channel := r.Notifier.Subscribe(postID)

	go func() {
		<-ctx.Done()
		logger.Info("Unsubscribing from comments for post", slog.String("func ", f), slog.Int("postID", postID))
		r.Notifier.Unsubscribe(postID, channel)
	}()

	return channel, nil
}
