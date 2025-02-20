package service

import (
	"sync"

	"github.com/riddion72/ozon_test/internal/domain"
)

type Notifier struct {
	mu          sync.RWMutex
	subscribers map[string][]chan domain.Comment
}

func NewNotifier() *Notifier {
	return &Notifier{
		subscribers: make(map[string][]chan domain.Comment),
	}
}

func (n *Notifier) Subscribe(postID string) chan domain.Comment {
	ch := make(chan domain.Comment, 1)

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscribers[postID] = append(n.subscribers[postID], ch)
	return ch
}

func (n *Notifier) Unsubscribe(postID string, ch chan domain.Comment) {
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

func (n *Notifier) Notify(postID string, comment domain.Comment) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, ch := range n.subscribers[postID] {
		ch <- comment
	}
}
