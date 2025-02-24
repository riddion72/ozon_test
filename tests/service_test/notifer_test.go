package service_test

import (
	"testing"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNotifier(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		n := service.NewNotifier()
		postID := 1

		ch := n.Subscribe(postID)

		comment := &domain.Comment{ID: 1}
		n.Notify(postID, comment)

		select {
		case msg := <-ch:
			require.Equal(t, comment, msg)
		case <-time.After(10 * time.Millisecond):
			require.Fail(t, "Timeout waiting for notification")
		}
	})

	t.Run("Unsubscribel", func(t *testing.T) {
		n := service.NewNotifier()
		postID := 1
		ch := n.Subscribe(postID)

		n.Unsubscribe(postID, ch)

		comment := &domain.Comment{ID: 1}
		n.Notify(postID, comment)

		select {
		case <-ch:
			require.Fail(t, "Should not receive message")
		case <-time.After(50 * time.Millisecond):
		}
	})

	t.Run("Notify", func(t *testing.T) {
		n := service.NewNotifier()
		postID := 1
		comment := &domain.Comment{ID: 1}

		ch1 := n.Subscribe(postID)
		ch2 := n.Subscribe(postID)

		n.Notify(postID, comment)

		require.Equal(t, comment, <-ch1)
		require.Equal(t, comment, <-ch2)
	})
}
