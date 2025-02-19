package subscriptions

import (
	"sync"

	"github.com/riddion72/ozon_test/internal/domain"
)

type SubscriptionManager struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *domain.Comment
}

func (sm *SubscriptionManager) Publish(postID string, comment *domain.Comment) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, ch := range sm.subscribers[postID] {
		ch <- comment
	}
}
