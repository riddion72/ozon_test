package app

import (
	"github.com/riddion72/ozon_test/internal/app/gqlserver"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/graph/resolvers"
	"github.com/riddion72/ozon_test/internal/service"
	"github.com/riddion72/ozon_test/internal/storage"
)

type App struct {
	Server *gqlserver.GQLServer
}

func NewApp(config *config.Config) *App {
	//Инициализация хранилища
	db := storage.NewStorage(config.DB)

	// Инициализация сервисов
	notifier := service.NewNotifier()
	postService := service.NewPostService(*db.Post)
	commentService := service.NewCommentService(*db.Comment, *db.Post)

	// Создаем резолвер
	resolver := resolvers.NewResolver(postService, commentService, notifier)

	// Cервер
	server := gqlserver.New(config.Server, resolver)

	return &App{Server: server}
}
