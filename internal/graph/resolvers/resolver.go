package resolvers

import (
	"github.com/riddion72/ozon_test/internal/storage"
)

type Resolver struct {
	PostStorage    storage.PostStorage
	CommentStorage storage.CommentStorage
	Subscriptions  *SubscriptionManager
}

// в resolvers.go
func (r *Resolver) Complexity() graphql.ComplexityRoot {
	return graphql.ComplexityRoot{
		// Пример: ограничить глубину вложенности комментариев
		Comment: func(childComplexity int) int {
			return childComplexity * 2 // Множитель сложности для вложенных комментариев
		},
	}
}
