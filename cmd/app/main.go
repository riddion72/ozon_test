package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/riddion72/ozon_test/internal/config"
	// "github.com/riddion72/ozon_test/internal/graph"
	"github.com/riddion72/ozon_test/internal/graph/generated"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/service"
	"github.com/riddion72/ozon_test/internal/storage"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()
	// Загрузка конфигурации
	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Инициализация логгера
	logger.MustInit(cfg.Logger.Level)
	logger.Info("Starting application", slog.String("version", "1.0.0"))

	//Инициализация хранилища
	postRepo, commentRepo := storage.CreateStorages(cfg.DB)

	// Инициализация сервисов
	notifier := service.NewNotifier()
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, notifier)

	// Создаем резолвер
	resolver := &graph.Resolver{
		Services: &service.Services{
			Posts:    postService,
			Comments: commentService,
			Notifier: notifier,
		},
	}

	// Настройка HTTP сервера
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers:  resolver,
		Complexity: graph.ComplexityConfig(),
	}))
	srv.Use(extension.FixedComplexityLimit(1000))

	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	httpServer := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", slog.String("address", cfg.Server.Address))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.String("error", err.Error()))
	}
	logger.Info("Server stopped")
}
