package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"

	"github.com/riddion72/ozon_test/internal/graph"
	"github.com/riddion72/ozon_test/internal/storage/inmem"
)

func main() {

	postStorage := inmem.NewPostStorage()
	resolver := &graph.Resolver{
		PostStorage: postStorage,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.Use(extension.FixedComplexityLimit(1000)) // Лимит сложности

	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
