package service

import (
	"sync"

	"github.com/riddion72/ozon_test/internal/domain"
)

type notifier struct {
	mu          sync.RWMutex
	subscribers map[int][]chan *domain.Comment
}

func NewNotifier() *notifier {
	return &notifier{
		subscribers: make(map[int][]chan *domain.Comment),
	}
}

func (n *notifier) Subscribe(postID int) chan *domain.Comment {
	ch := make(chan *domain.Comment, 1)

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscribers[postID] = append(n.subscribers[postID], ch)
	return ch
}

func (n *notifier) Unsubscribe(postID int, ch chan *domain.Comment) {
	n.mu.Lock()
	defer n.mu.Unlock()

	subscribers := n.subscribers[postID]
	for i, sub := range subscribers {
		if sub == ch {
			close(ch)
			n.subscribers[postID] = append(subscribers[:i], subscribers[i+1:]...)
			return
		}
	}
}

func (n *notifier) Notify(postID int, comment *domain.Comment) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, ch := range n.subscribers[postID] {
		ch <- comment
	}
}
