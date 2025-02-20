package resolvers

import (
	"github.com/riddion72/ozon_test/internal/service"
)

// Resolver корневая структура для всех резолверов
type Resolver struct {
	services *service.Services
}

// NewResolver конструктор для Resolver
func NewResolver(services *service.Services) *Resolver {
	return &Resolver{services: services}
}

func (r *Resolver) Complexity() graphql.ComplexityRoot {
	return graphql.ComplexityRoot{
		// ограничить глубину вложенности комментариев
		Comment: func(childComplexity int) int {
			return childComplexity * 2 // Множитель сложности для вложенных комментариев
		},
	}
}
