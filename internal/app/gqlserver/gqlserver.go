package gqlserver

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/graph/complexity"
	"github.com/riddion72/ozon_test/internal/graph/generated"
	"github.com/riddion72/ozon_test/internal/graph/resolvers"
	"github.com/riddion72/ozon_test/internal/logger"
)

type GQLServer struct {
	cfg        config.Server
	srv        *handler.Server
	httpServer *http.Server
}

func New(config config.Server, resolver *resolvers.Resolver) *GQLServer {

	cmplx := complexity.NewComplexity

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
		Complexity: generated.ComplexityRoot{
			Comment:      cmplx().Comment,
			Mutation:     cmplx().Mutation,
			Post:         cmplx().Post,
			Query:        cmplx().Query,
			Subscription: cmplx().Subscription,
		},
	}))
	srv.Use(extension.FixedComplexityLimit(config.ComplexityLimit))

	srv.AddTransport(&transport.Websocket{})

	return &GQLServer{
		cfg: config,
		srv: srv,
	}
}

func (a *GQLServer) run() error {
	const f = "gqlserver.run"
	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	router.Handle("/graphql", a.srv)
	router.Handle("/subscriptions", a.srv)

	a.httpServer = &http.Server{
		Addr:         a.cfg.Address,
		Handler:      router,
		ReadTimeout:  a.cfg.Timeout,
		WriteTimeout: a.cfg.Timeout,
		IdleTimeout:  a.cfg.IddleTimeout,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", slog.String("func", f), slog.String("address", a.cfg.Address))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.String("func", f), slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("GraphQL server started", slog.String("func", f), slog.String("port", a.cfg.Address))

	<-done
	logger.Info("Shutting down server...")

	if err := a.Shutdown(); err != nil {
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")
	return nil
}

func (a *GQLServer) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *GQLServer) Shutdown() error {
	const f = "gqlserver.Shutdown"
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	logger.Info("Stopping GraphQL server", slog.String("func", f))
	if err := a.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.String("func", f), slog.String("error", err.Error()))
		return err
	}
	return nil
}
