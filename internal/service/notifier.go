package service

import (
	"log"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
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
	ch := make(chan *domain.Comment, 1)

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscribers[postID] = append(n.subscribers[postID], ch)
	return ch
}

func (n *Notifier) Unsubscribe(postID int, ch chan *domain.Comment) {
	n.mu.Lock()
	defer n.mu.Unlock()

	subscribers := n.subscribers[postID]
	for i, sub := range subscribers {
		if sub == ch {
			n.subscribers[postID] = append(subscribers[:i], subscribers[i+1:]...)
			return
		}
	}
}

func (n *Notifier) Notify(postID int, comment *domain.Comment) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	subscribers := make([]chan *domain.Comment, len(n.subscribers[postID]))
	copy(subscribers, n.subscribers[postID])

	for _, ch := range subscribers {
		go func(c chan *domain.Comment) {
			select {
			case c <- comment:
			case <-time.After(100 * time.Millisecond):
				log.Println("Failed to deliver message")
			}
		}(ch)
	}
}
