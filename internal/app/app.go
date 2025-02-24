package app

import (
	"log/slog"

	"github.com/riddion72/ozon_test/internal/app/gqlserver"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/graph/resolvers"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/service"
	"github.com/riddion72/ozon_test/internal/storage"
)

type App struct {
	Server *gqlserver.GQLServer
}

func NewApp(config *config.Config) *App {
	const f = "NewApp"
	db := storage.NewStorage(config.DB)
	logger.Info("Database initialized successfully", slog.String("func", f))
	notifier := service.NewNotifier()
	postService := service.NewPostService(*db.Post)
	commentService := service.NewCommentService(*db.Comment, *db.Post)

	resolver := resolvers.NewResolver(postService, commentService, notifier)

	server := gqlserver.NewServer(config.Server, resolver)
	logger.Info("GraphQL server initialized successfully", slog.String("func", f))
	return &App{Server: server}
}
