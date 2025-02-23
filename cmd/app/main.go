package main

import (
	"flag"
	"log"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/app"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/logger"
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

	// init app
	app := app.NewApp(cfg)

	// run app
	app.Server.MustRun()

	// //Инициализация хранилища
	// postRepo, commentRepo := storage.CreateStorages(cfg.DB)

	// // Инициализация сервисов
	// notifier := service.NewNotifier()
	// postService := service.NewPostService(postRepo)
	// commentService := service.NewCommentService(commentRepo, postRepo)

	// // Создаем резолвер
	// resolver := resolvers.NewResolver(postService, commentService, notifier)
	// cmplx := complexity.NewComplexity

	// Настройка HTTP сервера
	// srv := handler.New(generated.NewExecutableSchema(generated.Config{
	// 	Resolvers: resolver,
	// 	Complexity: generated.ComplexityRoot{
	// 		Comment:      cmplx().Comment,
	// 		Mutation:     cmplx().Mutation,
	// 		Post:         cmplx().Post,
	// 		Query:        cmplx().Query,
	// 		Subscription: cmplx().Subscription,
	// 	},
	// }))
	// srv.Use(extension.FixedComplexityLimit(1000))

	// router := http.NewServeMux()
	// router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// router.Handle("/query", srv)

	// httpServer := &http.Server{
	// 	Addr:         cfg.Server.Address,
	// 	Handler:      router,
	// 	ReadTimeout:  cfg.Server.Timeout,
	// 	WriteTimeout: cfg.Server.Timeout,
	// 	IdleTimeout:  cfg.Server.IdleTimeout,
	// }

	// // Graceful shutdown
	// done := make(chan os.Signal, 1)
	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	logger.Info("Starting server", slog.String("address", cfg.Server.Address))
	// 	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		logger.Error("Failed to start server", slog.String("error", err.Error()))
	// 		os.Exit(1)
	// 	}
	// }()

	// <-done
	// logger.Info("Shutting down server...")

	// ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	// defer cancel()

	// if err := httpServer.Shutdown(ctx); err != nil {
	// 	logger.Error("Server shutdown failed", slog.String("error", err.Error()))
	// }
	// logger.Info("Server stopped")
}
