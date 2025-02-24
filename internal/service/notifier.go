package service

import (
	"log/slog"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
)

type Notifier struct {
	mu          sync.RWMutex
	subscribers map[int][]chan *domain.Comment
}

func NewNotifier() *Notifier {
	return &Notifier{
		subscribers: make(map[int][]chan *domain.Comment),
	}
}

func (n *Notifier) Subscribe(postID int) chan *domain.Comment {
	const f = "service.Notifier.Subscribe"
	ch := make(chan *domain.Comment, 1)

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscribers[postID] = append(n.subscribers[postID], ch)
	logger.Info("Subscribed to post", slog.String("func", f), slog.Int("postID", postID))
	return ch
}

func (n *Notifier) Unsubscribe(postID int, ch chan *domain.Comment) {
	const f = "service.Notifier.Unsubscribe"
	n.mu.Lock()
	defer n.mu.Unlock()

	subscribers := n.subscribers[postID]
	for i, sub := range subscribers {
		if sub == ch {
			n.subscribers[postID] = append(subscribers[:i], subscribers[i+1:]...)
			logger.Info("Unsubscribed from post", slog.String("func", f), slog.Int("postID", postID))
			return
		}
	}
}

func (n *Notifier) Notify(postID int, comment *domain.Comment) {
	const f = "service.Notifier.Notify"
	n.mu.RLock()
	defer n.mu.RUnlock()

	subscribers := make([]chan *domain.Comment, len(n.subscribers[postID]))
	copy(subscribers, n.subscribers[postID])

	for _, ch := range subscribers {
		go func(c chan *domain.Comment) {
			select {
			case c <- comment:
				logger.Info("Delivered comment to post", slog.String("func", f), slog.Int("postID", postID), slog.Int("commentID", comment.ID))
			case <-time.After(100 * time.Millisecond):
				logger.Warn("Failed to deliver message", slog.String("func", f), slog.Int("postID", postID))
			}
		}(ch)
	}
}
