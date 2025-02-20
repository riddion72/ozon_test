package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage/inmemory"
	"github.com/riddion72/ozon_test/internal/storage/postgres"
)

type PostStorage interface {
	Create(post domain.Post) error
	GetByID(id string) (domain.Post, bool)
	List(limit, offset int) []domain.Post
}

type CommentStorage interface {
	Create(comment domain.Comment) error
	GetByID(id string) (domain.Comment, bool)
	GetByPostID(postID string, limit, offset int) []domain.Comment
}

type Storage struct {
	Post    *PostStorage
	Comment *CommentStorage
}

func NewStorage(post *PostStorage, comment *CommentStorage) *Storage {
	return &Storage{Post: post, Comment: comment}
}

func CreateStorages(cfg config.DB) (PostStorage, CommentStorage) {
	var postRepo PostStorage
	var commentRepo CommentStorage
	var db *sqlx.DB

	// Выбор хранилища
	if cfg.Host != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// Подключение к PostgreSQL с повторами
		var err error
		db, err = postgres.ConnectWithRetries(ctx, cfg)
		if err != nil {
			logger.Error("Failed to connect to PostgreSQL: ", slog.String("error", err.Error()))
			// Если к PostgreSQL подключиться не получилось используем In-memory реализацию
			postRepo = inmemory.NewPostRepo()
			commentRepo = inmemory.NewCommentRepo()
			logger.Info("Using in-memory storage")
		} else {
			defer db.Close()

			postRepo = postgres.NewPostRepository(db)
			commentRepo = postgres.NewCommentRepository(db)
		}
	} else {
		// In-memory реализация
		postRepo = inmemory.NewPostRepo()
		commentRepo = inmemory.NewCommentRepo()
		logger.Info("Using in-memory storage")
	}

	return postRepo, commentRepo
}
